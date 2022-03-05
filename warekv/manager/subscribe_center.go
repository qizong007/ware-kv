package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	defaultDefaultCallbackMethod     = http.MethodGet
	defaultCallbackRetryQueueLen     = 128
	defaultCallbackRetryTickInterval = time.Second
)

// status
const (
	callbackCreated = iota // 回调任务创建（初始状态）
	callbackRequest        // 回调请求中
	callbackRetry          // 回调重试中
	callbackSuccess        // 回调成功（终态）
	callbackFail           // 回调失败（终态）
)

// event
const (
	CallbackSetEvent = iota
	CallbackDeleteEvent
)

var (
	center *SubscribeCenter
)

// SubscribeCenter 订阅中心
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

func NewSubscribeCenter(option *SubscribeCenterOption) *SubscribeCenter {
	defaultCallbackMethod := defaultDefaultCallbackMethod
	retryQueueLen := defaultCallbackRetryQueueLen
	retryTickInterval := defaultCallbackRetryTickInterval
	if option != nil {
		defaultCallbackMethod = option.DefaultCallbackMethod
		retryQueueLen = int(option.RetryQueueLen)
		retryTickInterval = time.Duration(option.RetryTickInterval)
	}
	center = &SubscribeCenter{
		record:                make(map[string][]*CallbackPlan),
		retryQueue:            make(chan *CallbackPlan, retryQueueLen),
		retryTicker:           time.NewTicker(retryTickInterval),
		retryCloser:           make(chan bool),
		defaultCallbackMethod: defaultCallbackMethod,
	}
	return center
}

func (s *SubscribeCenter) Start() {
	go center.scheduledRetry()
	fmt.Println("Subscriber's Retry worker starts working...")
}

func (s *SubscribeCenter) Close() {
	s.retryCloser <- true
}

func GetSubscribeCenter() *SubscribeCenter {
	return center
}

type SubscribeOption struct {
	Key          string
	CallbackPath string
	ExpectEvent  []int
	RetryTimes   int
}

// Subscribe 订阅
func (s *SubscribeCenter) Subscribe(option *SubscribeOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	plan := s.generateCallbackPlan(option.CallbackPath, ExpectEvent(option.ExpectEvent), RetryTimes(option.RetryTimes))
	key := option.Key
	if plans, ok := s.record[key]; ok {
		plans = append(plans, plan)
	} else {
		plans = []*CallbackPlan{plan}
		s.record[key] = plans
	}
}

// Notify 通知回调
func (s *SubscribeCenter) Notify(key string, newVal interface{}, event int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	plans := s.record[key]
	for i := range plans {
		if plans[i].status == callbackSuccess || plans[i].status == callbackFail {
			plans = append(plans[:i], plans[i+1:]...) // delete
		}
		if plans[i].status == callbackCreated {
			if isEventInList(event, *plans[i].expectEvent) {
				plans[i].notify(newVal)
			} else {
				plans = append(plans[:i], plans[i+1:]...) // delete
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

// 回调计划生成
func (s *SubscribeCenter) generateCallbackPlan(path string, options ...CallbackPlanOption) *CallbackPlan {
	plan := &CallbackPlan{
		center:         s,
		callbackPath:   path,
		callbackMethod: s.defaultCallbackMethod,
		status:         callbackCreated,
		expectEvent:    &[]int{CallbackSetEvent, CallbackDeleteEvent},
	}
	for _, option := range options {
		option(plan)
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

// CallbackPlan 回调计划
// 回调成功段时间后，统一卸载任务
type CallbackPlan struct {
	center         *SubscribeCenter
	callbackPath   string
	callbackMethod string
	param          interface{} // 参数缓存
	status         int
	expectEvent    *[]int
	retryTimes     int // 重试次数
	//expectValue
}

func (p *CallbackPlan) notify(newVal interface{}) {
	p.param = newVal
	p.status = callbackRequest
	// 按照回调方法分发请求
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
	// fixme: 解析newVal
	return fmt.Sprintf("%s?newVal=%s", p.callbackPath, paramStr), nil
}

func (p *CallbackPlan) notifyInGet() {
	getPath, err := p.generateGetPath()
	if err != nil {
		log.Println("generateGetPath(), json解析错误", err)
		return
	}
	resp, err := http.Get(getPath)
	if err != nil {
		p.dealWithCallbackErr()
		log.Println("http.Get(getPath), GET回调错误", err)
		return
	}
	defer resp.Body.Close()
	p.retryOrAbort(resp)
}

func (p *CallbackPlan) notifyWithBody() {
	req, err := p.generateRequest()
	if err != nil {
		log.Println("p.generateRequest(), 生成 req 错误", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.dealWithCallbackErr()
		log.Println("client.Do(req) 回调错误", err)
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
	if p.retryTimes > 0 {
		p.retryTimes--
		p.status = callbackRetry
		p.center.retryQueue <- p
	} else {
		p.status = callbackFail
	}
}

type CallbackPlanOption func(*CallbackPlan)

func ExpectEvent(events []int) CallbackPlanOption {
	return func(p *CallbackPlan) {
		if events != nil {
			*p.expectEvent = events
		}
	}
}

func RetryTimes(times int) CallbackPlanOption {
	return func(p *CallbackPlan) {
		if times != 0 {
			p.retryTimes = times
		}
	}
}
