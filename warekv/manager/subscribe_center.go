package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"ware-kv/warekv/util"
)

const (
	defaultDefaultCallbackMethod     = http.MethodGet
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

func (s *SubscribeCenter) start() {
	go center.scheduledRetry()
	fmt.Println("Subscriber's Retry worker starts working...")
}

func (s *SubscribeCenter) Close() {
	s.retryCloser <- true
}

type SubscribeManifest struct {
	Key          string `json:"k"`
	CallbackPath string `json:"cp"`
	ExpectEvent  []int  `json:"e"`
	RetryTimes   int    `json:"rt"`
	IsPersistent bool   `json:"ip"`
}

func (s *SubscribeCenter) Subscribe(option *SubscribeManifest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cpOption := &CallbackPlanOption{
		callbackPath: option.CallbackPath,
		events:       option.ExpectEvent,
		retryTimes:   option.RetryTimes,
		isPersistent: option.IsPersistent,
	}
	plan := s.generateCallbackPlan(cpOption)
	key := option.Key
	if plans, ok := s.record[key]; ok {
		plans = append(plans, plan)
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
				deleteCallbackPlan(plans, i)
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
		callbackPath:   option.callbackPath,
		callbackMethod: s.defaultCallbackMethod,
		status:         callbackCreated,
		expectEvent:    &[]int{CallbackSetEvent, CallbackDeleteEvent},
		isPersistent:   option.isPersistent,
	}
	if option.retryTimes != 0 {
		plan.retryTimes = option.retryTimes
		plan.leftRetryTimes = option.retryTimes
	}
	if option.events != nil {
		*plan.expectEvent = option.events
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
	//expectValue
}

func (p *CallbackPlan) notify(newVal interface{}) {
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
	param, err := json.Marshal(p.param)
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
	// fixme: resolve newVal
	return fmt.Sprintf("%s?newVal=%s", p.callbackPath, paramStr), nil
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
	callbackPath string
	events       []int
	retryTimes   int
	isPersistent bool
}
