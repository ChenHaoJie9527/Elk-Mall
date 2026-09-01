package response

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
