package response

import (
	"net/http"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/labstack/echo/v5"
)

type Response struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
	ErrMsg string `json:"err_msg,omitempty"`
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

// WriteResponse 统一写入业务 JSON。e 为 nil 时按成功（errno.OK）处理。
func WriteResponse(c *echo.Context, data any, e *errno.Errno) error {
	if e == nil {
		e = errno.OK
	}
	return c.JSON(http.StatusOK, Response{
		Code:   e.Code,
		Msg:    e.Msg,
		ErrMsg: e.ErrMsg,
		Data:   data,
	})
}
