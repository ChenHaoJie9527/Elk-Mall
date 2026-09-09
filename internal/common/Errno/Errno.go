package errno

// 定义 Errno 类型
type Errno struct {
	Code   int
	Msg    string
	ErrMsg string
}

func (e *Errno) Error() string { return e.Msg }

func (e *Errno) clone() *Errno {
	cp := *e
	return &cp
}

// WithMsg 只改文案、保持原 Code，返回副本，不修改全局 Errno。
func (e *Errno) WithMsg(msg string) *Errno {
	cp := e.clone()
	cp.Msg = msg
	return cp
}

// WithErrMsg 把底层错误写进 ErrMsg，返回副本，不修改全局 Errno。
func (e *Errno) WithErrMsg(rawErr error) *Errno {
	cp := e.clone()
	if rawErr != nil {
		cp.ErrMsg = rawErr.Error()
	}
	return cp
}

// IsOK 判断 Errno 是否是 OK
func (e *Errno) IsOK() bool { return e.Code == 200 }

// 定义常用的 Errno
var (
	OK                      = &Errno{Code: 0, Msg: "OK"}
	ParmError               = &Errno{Code: 10001, Msg: "参数错误"}
	NotFound                = &Errno{Code: 10002, Msg: "资源不存在"}
	ServerError             = &Errno{Code: 10003, Msg: "服务器内部错误"}
	UsernameExists          = &Errno{Code: 10004, Msg: "用户名已存在"}
	UsernameOrPasswordError = &Errno{Code: 10005, Msg: "用户名或密码错误"}
	UserNotFound            = &Errno{Code: 10006, Msg: "用户不存在"}
	TokenExpired            = &Errno{Code: 10007, Msg: "Token 已过期"}
)
