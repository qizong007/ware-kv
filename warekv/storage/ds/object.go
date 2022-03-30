package ds

import (
	"github.com/qizong007/ware-kv/warekv/util"
	"sync"
)

type Object struct {
	Base
	object map[string]interface{}
	rw     sync.RWMutex
}

func (o *Object) GetValue() interface{} {
	o.rw.RLock()
	defer o.rw.RUnlock()
	return o.object
}

func MakeObject(object map[string]interface{}) *Object {
	return &Object{
		Base:   *NewBase(util.ObjectDS),
		object: object,
	}
}

func (o *Object) GetFieldByKey(key string) interface{} {
	o.rw.RLock()
	defer o.rw.RUnlock()
	return o.object[key]
}

func (o *Object) SetFieldByKey(key string, val interface{}) {
	o.rw.Lock()
	defer o.rw.Unlock()
	o.Update()
	o.object[key] = val
}
