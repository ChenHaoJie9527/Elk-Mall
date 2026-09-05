package dto

// 注册请求 请求参数校验
type RegisterReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// 登录请求 请求参数校验
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 登录响应 响应参数
type LoginResp struct {
	Token string   `json:"token"`
	User  UserResp `json:"user"`
}

// 用户响应 响应参数
type UserResp struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
