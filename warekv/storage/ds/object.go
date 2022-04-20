package ds

import (
	"github.com/qizong007/ware-kv/warekv/util"
	"sync"
	"unsafe"
)

type Object struct {
	Base
	object map[string]interface{}
	rw     sync.RWMutex
}

var objectStructMemUsage int

func init() {
	objectStructMemUsage = int(unsafe.Sizeof(Object{}))
}

func (o *Object) GetValue() interface{} {
	o.rw.RLock()
	defer o.rw.RUnlock()
	return o.object
}

func (o *Object) Size() int {
	size := objectStructMemUsage
	if o.ExpireTime != nil {
		size += 8
	}
	o.rw.RLock()
	defer o.rw.RUnlock()
	if rSize := util.GetRealSizeOf(o.object); rSize > 0 {
		size += rSize
	}
	return size
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

func (o *Object) DeleteFieldByKey(key string) {
	o.rw.Lock()
	defer o.rw.Unlock()
	o.Update()
	delete(o.object, key)
}
