package manager

const (
	defaultCallbackMethod = "POST"
)

// status
const (
	callbackCreated = iota // 回调任务创建（初始状态）
	callbackReady          // 回调就绪
	callbackRequest        // 回调请求中
	callbackRetry          // 回调重试中
	callbackSuccess        // 回调成功（终态）
	callbackFail           // 回调失败（终态）
)

// event
const (
	callbackSet = iota
	callbackDelete
)

// SubscribeCenter 订阅中心
type SubscribeCenter struct {
	record map[string][]*CallbackPlan
}

// CallbackPlan 回调计划
// 回调成功段时间后，统一卸载任务
type CallbackPlan struct {
	callbackPath string
	params       map[string]interface{}
	status       int
	expectEvent  *[]int
	//expectValue
}

// Subscribe 订阅
func (s *SubscribeCenter) Subscribe(key string, plan *CallbackPlan) {

}
