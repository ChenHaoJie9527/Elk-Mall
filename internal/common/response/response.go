package response

import errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Success 成功响应
func Success(data any) *Response {
	return &Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

func Err(e *errno.Errno) *Response {
	return &Response{Code: e.Code, Msg: e.Msg}
}
