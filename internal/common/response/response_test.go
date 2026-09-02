package response

import (
	"encoding/json"
	"testing"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
)

// TestErr_MapsErrnoToResponse 测试 Err 方法是否将 Errno 映射为 Response
func TestErr_MapsErrnoToResponse(t *testing.T) {
	tests := []struct {
		name string
		in   *errno.Errno
		code int
		msg  string
	}{
		{"参数错误", errno.ParmError, 10001, "参数错误"},
		{"不存在", errno.NotFound, 10002, "资源不存在"},
		{"内部错误", errno.ServerError, 10003, "服务器内部错误"},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			got := Err(v.in)
			if got.Code != v.code {
				t.Errorf("expected code %d, got %d", v.code, got.Code)
			}
			if got.Msg != v.msg {
				t.Errorf("expected msg %s, got %s", v.msg, got.Msg)
			}
			if got.Data != nil {
				t.Errorf("expected data nil, got %v", got.Data)
			}
		})
	}
}

// TestErr_WithMsgKeepsCode 测试 WithMsg 方法是否保持了 Errno 的 Code
func TestErr_WithMsgKeepsCode(t *testing.T) {
	got := Err(errno.ParmError.WithMsg("id 必填"))
	if got.Code != errno.ParmError.Code {
		t.Errorf("expected code %d, got %d", errno.ParmError.Code, got.Code)
	}
	if got.Msg != "id 必填" {
		t.Errorf("expected msg %s, got %s", "id 必填", got.Msg)
	}
}

// TestErr_JSONShape 测试 Err 方法返回的 JSON 格式是否正确
func TestErr_JSONShape(t *testing.T) {
	// Marshal 将结构体转换为 JSON 字符串
	data, err := json.Marshal(Err(errno.NotFound))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"code":10002,"msg":"资源不存在","data":null}`
	if string(data) != want {
		t.Errorf("expected %s, got %s", want, string(data))
	}
}
