package middleware

import (
	"errors"
	"net/http"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/labstack/echo/v5"
)

// HTTPErrorHandler 自定义 HTTP 错误处理函数
// 1.已提交响应，则直接返回
// 2.业务错误走errno，返回json响应
// 3.其余走 http 状态
func HTTPErrorHandler(c *echo.Context, err error) {

	// 如果响应已发送，则直接返回
	if resp, _ := echo.UnwrapResponse(c.Response()); resp != nil && resp.Committed {
		return
	}

	var e *errno.Errno
	// 通过 as 断言，如果 err 是 errno.Errno 类型，则将 errno.Errno 赋值给 e
	if errors.As(err, &e) {
		_ = c.JSON(http.StatusOK, response.Err(e))
		return
	}

	// 通过 错误 获取状态码
	code := echo.StatusCode(err)
	if code == 0 {
		// 如果状态码为 0，则设置为 500
		code = http.StatusInternalServerError
	}
	_ = c.JSON(code, response.Response{
		Code: code,
		Msg:  http.StatusText(code),
	})

}
