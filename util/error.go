package util

const (
	Success       = 0
	ParamError    = 10000
	KeyNotExisted = 20000
)

var ErrCode2Msg = map[int]string{
	Success:       "success",
	ParamError:    "something wrong with param...",
	KeyNotExisted: "key is not existed!",
}
