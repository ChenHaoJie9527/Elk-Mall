package response

import (
	"errors"
	"net/http"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/labstack/echo/v5"
)

// HTTPErrorHandler 是 Echo 的全局错误处理函数，不是中间件。
// 1. 响应已提交则直接返回
// 2. 业务错误走 errno，返回 JSON
// 3. 其余走 HTTP 状态
func HTTPErrorHandler(c *echo.Context, err error) {
	if resp, _ := echo.UnwrapResponse(c.Response()); resp != nil && resp.Committed {
		return
	}

	var e *errno.Errno
	if errors.As(err, &e) {
		_ = c.JSON(http.StatusOK, Err(e))
		return
	}

	code := echo.StatusCode(err)
	if code == 0 {
		code = http.StatusInternalServerError
	}
	_ = c.JSON(code, Response{
		Code: code,
		Msg:  http.StatusText(code),
	})
}
