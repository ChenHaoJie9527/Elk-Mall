package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/labstack/echo/v5"
)

// 构造一次请求用的 Echo Context 和响应记录器
func newCtx() (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) response.Response {
	t.Helper()
	var body response.Response
	// 解析响应体并复制到 body 中
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal %q: %v", rec.Body.String(), err)
	}
	return body
}

// 测试已提交响应，则直接返回
func TestHTTPErrorHandler_CommittedResponseNotOverwritten(t *testing.T) {
	c, rec := newCtx()
	// 先提交一次
	if err := c.JSON(http.StatusOK, response.Success("already")); err != nil {
		t.Fatalf("seed json: %v", err)
	}
	HTTPErrorHandler(c, errno.ParmError)

	if rec.Code != http.StatusOK {
		t.Errorf("HTTP = %d, want 200", rec.Code)
	}
	body := decodeBody(t, rec)
	if body.Code != 0 || body.Msg != "success" || body.Data != "already" {
		t.Errorf("handler 覆盖了已提交响应: %+v", body)
	}
}

// 测试业务错误走errno，返回json响应
func TestHTTPErrorHandler_Errno(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code int
		msg  string
	}{
		{"参数错误", errno.ParmError, 10001, "参数错误"},
		{"资源不存在", errno.NotFound, 10002, "资源不存在"},
		{"服务器内部错误", errno.ServerError, 10003, "服务器内部错误"},
		{"WithMsg 只改文案", errno.ParmError.WithMsg("id 必填"), 10001, "id 必填"},
		{"wrap 后仍能识别", fmt.Errorf("bind: %w", errno.NotFound), 10002, "资源不存在"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newCtx()
			HTTPErrorHandler(c, tt.err)

			if rec.Code != http.StatusOK {
				t.Errorf("HTTP = %d, want 200", rec.Code)
			}
			body := decodeBody(t, rec)
			if body.Code != tt.code {
				t.Errorf("code = %d, want %d", body.Code, tt.code)
			}
			if body.Msg != tt.msg {
				t.Errorf("msg = %q, want %q", body.Msg, tt.msg)
			}
			if body.Data != nil {
				t.Errorf("data = %v, want nil", body.Data)
			}
		})
	}
}

// 测试其他错误走 http 状态码，返回json响应
func TestHTTPErrorHandler_HTTPStatus(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		msg    string
	}{
		{"echo 404", echo.ErrNotFound, http.StatusNotFound, http.StatusText(http.StatusNotFound)},
		{"NewHTTPError 400", echo.NewHTTPError(http.StatusBadRequest, "ignored"), http.StatusBadRequest, http.StatusText(http.StatusBadRequest)},
		{"未知 error 打成 500", errors.New("disk full"), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newCtx()
			HTTPErrorHandler(c, tt.err)

			if rec.Code != tt.status {
				t.Errorf("HTTP = %d, want %d", rec.Code, tt.status)
			}
			body := decodeBody(t, rec)
			if body.Code != tt.status {
				t.Errorf("code = %d, want %d", body.Code, tt.status)
			}
			if body.Msg != tt.msg {
				t.Errorf("msg = %q, want %q", body.Msg, tt.msg)
			}
		})
	}
}
