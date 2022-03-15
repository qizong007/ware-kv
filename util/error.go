package util

const (
	Success            = 0
	ParamError         = 10000
	KeyNotExisted      = 20000
	KeyHasDeleted      = 20001
	KeyHasExpired      = 20002
	ValueTypeError     = 30001
	TypeTransformError = 70001
	ScopeError         = 80001
)

var ErrCode2Msg = map[int]string{
	Success:            "success",
	ParamError:         "something wrong with param...",
	KeyNotExisted:      "key is not existed!",
	KeyHasDeleted:      "key has been deleted!",
	KeyHasExpired:      "key has been expired!",
	ValueTypeError:     "Error request type!",
	TypeTransformError: "something wrong with type transform...",
	ScopeError:         "scope is not in a right way!",
}
