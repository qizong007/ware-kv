package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"log"
	"time"
	"ware-kv/tracker"
	"ware-kv/util"
	"ware-kv/warekv/ds"
)

const (
	defaultTimeLimit     = 15 // time.Second
	defaultRetryTimes    = 2
	retryTimesLimit      = 15
	defaultRetryInterval = 1000 // time.Millisecond
	retryIntervalLimit   = 3000 // time.Millisecond
)

type LockParam struct {
	TimeLimit     int64  `json:"t" binding:"-"`
	ExpireTime    int64  `json:"expire_time" binding:"-"`
	RetryTimes    uint64 `json:"retry_times" binding:"-"`
	RetryInterval uint64 `json:"retry_interval" binding:"-"`
}

type UnlockParam struct {
	Guid string `json:"guid"`
}

func Lock(c *gin.Context) {
	param := LockParam{}
	_ = c.ShouldBindJSON(&param)

	var lock *ds.Lock

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}

	guid := xid.New().String()
	timeLimit := int64(defaultTimeLimit)
	retryTimes := uint64(defaultRetryTimes)
	retryInterval := uint64(defaultRetryInterval)
	if param.TimeLimit != 0 {
		timeLimit = param.TimeLimit
	}
	if param.RetryTimes != 0 {
		retryTimes = param.RetryTimes
		if retryTimes > retryTimesLimit {
			retryTimes = retryTimesLimit
		}
	}
	if param.RetryInterval != 0 {
		retryInterval = param.RetryInterval
		if retryInterval > retryIntervalLimit {
			retryInterval = retryIntervalLimit
		}
	}

	if isValNil(val) {
		// new
		lock = ds.MakeLock()
		cmd := tracker.NewCreateCommand(key.GetKey(), tracker.LockStruct, nil, lock.CreateTime, param.ExpireTime)
		setInTime(key, lock, param.ExpireTime, cmd)
	} else {
		// update
		if !isKVEffective(c, val) {
			return
		}
		if !isKVTypeCorrect(c, val, ds.LockDS) {
			return
		}
		lock = ds.Value2Lock(val)
	}
	err = lock.Lock(timeLimit, guid)
	if err != nil {
		err = retryLock(lock, timeLimit, guid, retryTimes, retryInterval)
		if err != nil {
			util.MakeResponse(c, &util.WareResponse{
				Code: util.LockRaceError,
			})
			return
		}
	}
	setNotify(key, lock)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  guid,
	})
}

func Unlock(c *gin.Context) {
	param := UnlockParam{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, ds.LockDS) {
		return
	}

	lock := ds.Value2Lock(val)

	err = lock.Unlock(param.Guid)
	if err != nil {
		util.MakeResponse(c, &util.WareResponse{
			Code: util.LockReleaseError,
		})
		return
	}
	setNotify(key, lock)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func retryLock(lock *ds.Lock, limit int64, guid string, times uint64, interval uint64) error {
	if times <= 0 {
		return fmt.Errorf("LockRaceFailed")
	}
	for times > 0 {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		err := lock.Lock(limit, guid)
		if err == nil {
			return nil
		}
		times--
	}
	return fmt.Errorf("LockRaceFailed")
}
