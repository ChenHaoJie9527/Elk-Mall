package errno

// 定义 Errno 类型
type Errno struct {
	Code int
	Msg  string
}

func (e *Errno) Error() string { return e.Msg }

// WithMsg 返回一个带有自定义消息的 Errno
func (e *Errno) WithMsg(msg string) *Errno {
	return &Errno{
		Code: e.Code,
		Msg:  msg,
	}
}

// 定义常用的 Errno
var (
	OK                      = &Errno{Code: 0, Msg: "OK"}
	ParmError               = &Errno{Code: 10001, Msg: "参数错误"}
	NotFound                = &Errno{Code: 10002, Msg: "资源不存在"}
	ServerError             = &Errno{Code: 10003, Msg: "服务器内部错误"}
	UsernameExists          = &Errno{Code: 10004, Msg: "用户名已存在"}
	UsernameOrPasswordError = &Errno{Code: 10005, Msg: "用户名或密码错误"}
	UserNotFound            = &Errno{Code: 10006, Msg: "用户不存在"}
)
