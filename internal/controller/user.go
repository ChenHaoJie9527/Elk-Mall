package controller

import (
	"net/http"
	"strconv"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
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
// 请求参数: username, password
// 返回参数: user
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
		return err
	}

	return c.JSON(http.StatusOK, response.Success(resp))
}

// 登录用户
// 请求参数: username, password
// 返回参数: token, user
func (u *User) Login(c *echo.Context) error {
	var req dto.LoginReq
	// bind: 绑定请求体到 req 中
	if err := c.Bind(&req); err != nil {
		return err
	}

	// 校验 参数: 使用 validator 校验参数
	// TODO: 后续 validate tag 校验参数
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, errno.ParmError)
	}

	resp, err := u.Svc.Login(&req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success(resp))
}

// 根据ID获取用户
// 请求参数: id
// 返回参数: user
func (u *User) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, errno.ParmError)
	}

	// 将字符串转换为 uint
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errno.ParmError)
	}

	// 获取用户信息
	resp, err := u.Svc.GetByID(uint(idUint))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.Success(resp))
}
