package manager

import (
	"ware-kv/warekv/storage"
)

// SubscribeCenter 订阅中心
type SubscribeCenter struct {
	// key -> []*CallbackPlan
}

// CallbackPlan 回调计划
// 回调成功段时间后，统一卸载任务
type CallbackPlan struct {
	// 回调地址(http）
	// 回调参数
	// 是否成功回调标志
	// 失败重试次数(optional)
	// 期望值(optional)
}

// Subscribe 订阅
func (s *SubscribeCenter) Subscribe(k *storage.Key, plan *CallbackPlan) {

}
