package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/util"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const (
	defaultDefaultCallbackMethod     = http.MethodPost
	defaultCallbackRetryQueueLen     = 128
	callbackRetryQueueLenMin         = 1
	callbackRetryQueueLenMax         = 64 * 1024 * 1024
	defaultCallbackRetryTickInterval = 1000
	callbackRetryTickIntervalMin     = 200
	callbackRetryTickIntervalMax     = 5000
)

// status
const (
	callbackCreated = iota // callback task created (Start)
	callbackRequest        // callback task requesting
	callbackRetry          // callback task retrying
	callbackSuccess        // callback task SUCCESS (Final)
	callbackFail           // callback task FAIL (Final)
)

// event
const (
	CallbackSetEvent = iota
	CallbackDeleteEvent
)

var (
	center *SubscribeCenter
)

type SubscribeCenter struct {
	record                map[string][]*CallbackPlan
	mu                    sync.Mutex
	retryQueue            chan *CallbackPlan
	retryTicker           *time.Ticker
	retryCloser           chan bool
	defaultCallbackMethod string
}

type SubscribeCenterOption struct {
	DefaultCallbackMethod string `yaml:"DefaultCallbackMethod"`
	RetryQueueLen         uint   `yaml:"RetryQueueLen"`
	RetryTickInterval     uint   `yaml:"RetryTickInterval"`
}

func DefaultSubscribeCenterOption() *SubscribeCenterOption {
	return &SubscribeCenterOption{
		DefaultCallbackMethod: defaultDefaultCallbackMethod,
		RetryQueueLen:         defaultCallbackRetryQueueLen,
		RetryTickInterval:     defaultCallbackRetryTickInterval,
	}
}

func NewSubscribeCenter(option *SubscribeCenterOption) *SubscribeCenter {
	defaultCallbackMethod := defaultDefaultCallbackMethod
	retryQueueLen := defaultCallbackRetryQueueLen
	retryTickInterval := time.Millisecond * time.Duration(defaultCallbackRetryTickInterval)
	if option != nil {
		defaultCallbackMethod = option.DefaultCallbackMethod
		retryQueueLen = util.SetIfHitLimit(int(option.RetryQueueLen), callbackRetryQueueLenMin, callbackRetryQueueLenMax)
		tickInterval := util.SetIfHitLimit(int(option.RetryTickInterval), callbackRetryTickIntervalMin, callbackRetryTickIntervalMax)
		retryTickInterval = time.Millisecond * time.Duration(tickInterval)
	}
	center = &SubscribeCenter{
		record:                make(map[string][]*CallbackPlan),
		retryQueue:            make(chan *CallbackPlan, retryQueueLen),
		retryTicker:           time.NewTicker(retryTickInterval),
		retryCloser:           make(chan bool),
		defaultCallbackMethod: defaultCallbackMethod,
	}
	center.start()
	return center
}

func GetSubscribeCenter() *SubscribeCenter {
	return center
}

func (s *SubscribeCenter) DefaultCallbackMethod() string {
	return s.defaultCallbackMethod
}

func (s *SubscribeCenter) start() {
	go center.scheduledRetry()
	log.Println("Subscriber's Retry worker starts working...")
}

func (s *SubscribeCenter) Close() {
	s.retryCloser <- true
}

type SubscribeManifest struct {
	Key          string      `json:"k"`
	CallbackPath string      `json:"cp"`
	ExpectEvent  []int       `json:"ee"`
	RetryTimes   int         `json:"rt"`
	IsPersistent bool        `json:"ip"`
	Method       string      `json:"m"`
	ExpectValue  interface{} `json:"ev"`
}

func (s *SubscribeCenter) Subscribe(option *SubscribeManifest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cpOption := &CallbackPlanOption{
		CallbackPath: option.CallbackPath,
		Events:       option.ExpectEvent,
		RetryTimes:   option.RetryTimes,
		IsPersistent: option.IsPersistent,
		Method:       option.Method,
		ExpectValue:  option.ExpectValue,
	}
	plan := s.generateCallbackPlan(cpOption)
	key := option.Key
	if plans, ok := s.record[key]; ok {
		s.record[key] = append(s.record[key], plan)
	} else {
		plans = []*CallbackPlan{plan}
		s.record[key] = plans
	}
}

func deleteCallbackPlan(plans []*CallbackPlan, i int) {
	plans = append(plans[:i], plans[i+1:]...)
}

func refreshCallbackPlan(plan *CallbackPlan) {
	plan.param = nil
	plan.status = callbackCreated
	plan.leftRetryTimes = plan.retryTimes
}

// Notify just callback
func (s *SubscribeCenter) Notify(key string, newVal interface{}, event int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	plans := s.record[key]
	for i := range plans {
		if plans[i].status == callbackSuccess || plans[i].status == callbackFail {
			if plans[i].isPersistent {
				refreshCallbackPlan(plans[i])
			} else {
				deleteCallbackPlan(plans, i)
				continue
			}
		}
		if plans[i].status == callbackCreated {
			if isEventInList(event, *plans[i].expectEvent) {
				plans[i].notify(newVal)
			} else {
				plans[i].status = callbackFail
			}
		}
	}
}

func isEventInList(event int, list []int) bool {
	if list == nil {
		return false
	}
	for i := range list {
		if list[i] == event {
			return true
		}
	}
	return false
}

func (s *SubscribeCenter) generateCallbackPlan(option *CallbackPlanOption) *CallbackPlan {
	plan := &CallbackPlan{
		center:         s,
		callbackPath:   option.CallbackPath,
		callbackMethod: option.Method,
		status:         callbackCreated,
		expectEvent:    &[]int{CallbackSetEvent, CallbackDeleteEvent},
		isPersistent:   option.IsPersistent,
		expectValue:    option.ExpectValue,
	}
	if option.RetryTimes != 0 {
		plan.retryTimes = option.RetryTimes
		plan.leftRetryTimes = option.RetryTimes
	}
	if option.Events != nil {
		*plan.expectEvent = option.Events
	}
	return plan
}

func (s *SubscribeCenter) scheduledRetry() {
	for {
		select {
		case <-s.retryTicker.C:
			if len(s.retryQueue) == 0 {
				continue
			}
			s.mu.Lock()
			for plan := range s.retryQueue {
				plan.retry()
				if len(s.retryQueue) == 0 {
					break
				}
			}
			s.mu.Unlock()
		case <-s.retryCloser:
			log.Println("Subscriber's Retry worker stops working...")
			return
		}
	}
}

// CallbackPlan
// If the callback plan execute SUCCESS, the plan will be uninstalled
type CallbackPlan struct {
	center         *SubscribeCenter
	callbackPath   string
	callbackMethod string
	param          interface{} // param cache
	status         int
	expectEvent    *[]int
	retryTimes     int
	leftRetryTimes int
	isPersistent   bool
	expectValue    interface{}
}

func (p *CallbackPlan) notify(newVal interface{}) {
	if !reflect.DeepEqual(newVal, p.expectValue) {
		return
	}
	p.param = newVal
	p.status = callbackRequest
	// distribute by callback method
	p.distribute()
}

func (p *CallbackPlan) retry() {
	p.distribute()
}

func (p *CallbackPlan) distribute() {
	switch p.callbackMethod {
	case http.MethodGet:
		p.notifyInGet()
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		p.notifyWithBody()
	}
}

func (p *CallbackPlan) generateRequest() (*http.Request, error) {
	ctx := map[string]interface{}{
		"new_val": p.param,
	}
	param, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(param)
	req, err := http.NewRequest(p.callbackMethod, p.callbackPath, reader)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (p *CallbackPlan) generateGetPath() (string, error) {
	paramStr, err := json.Marshal(p.param)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?new_val=%s", p.callbackPath, paramStr), nil
}

func (p *CallbackPlan) notifyInGet() {
	getPath, err := p.generateGetPath()
	if err != nil {
		log.Println("generateGetPath(), json.Unmarshall Fail", err)
		return
	}
	resp, err := http.Get(getPath)
	if err != nil {
		p.dealWithCallbackErr()
		log.Println("http.Get(getPath), GET callback Fail", err)
		return
	}
	defer resp.Body.Close()
	p.retryOrAbort(resp)
}

func (p *CallbackPlan) notifyWithBody() {
	req, err := p.generateRequest()
	if err != nil {
		log.Println("p.generateRequest(), generate 'req' Fail", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.dealWithCallbackErr()
		log.Println("client.Do(req) callback FAIL", err)
		return
	}
	defer resp.Body.Close()
	p.retryOrAbort(resp)
}

func (p *CallbackPlan) retryOrAbort(resp *http.Response) {
	if resp.StatusCode == http.StatusOK {
		p.status = callbackSuccess
	} else {
		p.dealWithCallbackErr()
	}
}

func (p *CallbackPlan) dealWithCallbackErr() {
	if p.leftRetryTimes > 0 {
		p.leftRetryTimes--
		p.status = callbackRetry
		p.center.retryQueue <- p
	} else {
		p.status = callbackFail
	}
}

type CallbackPlanOption struct {
	CallbackPath string
	Events       []int
	RetryTimes   int
	IsPersistent bool
	Method       string
	ExpectValue  interface{}
}

func (s *SubscribeCenter) View() []byte {
	data := make([]byte, 0)
	// subscribe center flag
	data = append(data, uint8(storage.SubscribeCenterFlag))
	s.mu.Lock()
	defer s.mu.Unlock()
	// sub keys num
	data = append(data, util.IntToBytes(len(s.record))...)
	// sub kv pairs
	for key, callbackPlanList := range s.record {
		data = append(data, subKVPairView(key, callbackPlanList)...)
	}
	return data
}

func subKVPairView(key string, callbackPlans []*CallbackPlan) []byte {
	// key (key len bytes)
	keyBytes := []byte(key)

	// value (value len bytes)
	valueBytes, err := callbackPlanListView(callbackPlans)
	if err != nil {
		log.Println(key, "get callbackPlanListView failed!")
		return []byte{}
	}

	// key len (4 bytes)
	keyLen := len(keyBytes)

	// value len (4 bytes)
	valueLen := len(valueBytes)

	data := make([]byte, 0, 8+keyLen+valueLen)
	data = append(data, util.IntToBytes(keyLen)...)
	data = append(data, keyBytes...)
	data = append(data, util.IntToBytes(valueLen)...)
	data = append(data, valueBytes...)

	return data
}

func callbackPlanListView(callbackPlans []*CallbackPlan) ([]byte, error) {
	data := make([]byte, 0)
	for _, callbackPlan := range callbackPlans {
		cpo := &CallbackPlanOption{
			CallbackPath: callbackPlan.callbackPath,
			Events:       *callbackPlan.expectEvent,
			RetryTimes:   callbackPlan.leftRetryTimes,
			IsPersistent: callbackPlan.isPersistent,
			Method:       callbackPlan.callbackMethod,
		}
		cpoBytes, err := json.Marshal(cpo)
		if err != nil {
			return []byte{}, err
		}
		data = append(data, cpoBytes...)
	}
	return data, nil
}

func (s *SubscribeCenter) GetFlag() storage.Flag {
	return storage.SubscribeCenterFlag
}
