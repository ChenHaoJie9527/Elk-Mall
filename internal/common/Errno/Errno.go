package errno

// 定义 Errno 类型
type Errno struct {
	Code int
	Msg  string
}

func (e *Errno) Error() string { return e.Msg }

func (e *Errno) WithMsg(msg string) *Errno {
	return &Errno{
		Code: e.Code,
		Msg:  msg,
	}
}

var (
	OK          = &Errno{Code: 0, Msg: "OK"}
	ParmError   = &Errno{Code: 10001, Msg: "参数错误"}
	NotFound    = &Errno{Code: 10002, Msg: "资源不存在"}
	ServerError = &Errno{Code: 10003, Msg: "服务器内部错误"}
)
