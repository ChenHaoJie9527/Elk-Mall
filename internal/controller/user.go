package controller

import (
	"net/http"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/dto"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/service"
	"github.com/labstack/echo/v5"
)

type User struct {
	Svc *service.UserService
}

func NewUserController(svc *service.UserService) *User {
	return &User{Svc: svc}
}

// 注册用户
func (u *User) Register(c *echo.Context) error {
	var req dto.RegisterReq
	// bind: 绑定请求体到 req 中
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errno.ParmError)
	}

	// 校验 参数: 使用 validator 校验参数
	// TODO: 后续 validate tag 校验参数
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, errno.ParmError)
	}

	resp, err := u.Svc.Register(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errno.ServerError)
	}

	return c.JSON(http.StatusOK, resp)
}
