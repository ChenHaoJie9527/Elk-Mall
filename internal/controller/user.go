package controller

import (
	"errors"
	"net/http"
	"strconv"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/dto"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/pkg"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/service"
	"github.com/golang-jwt/jwt/v5"
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

	// 从上下文中获取 token
	// “user”: 指的是 JWT 中间件中设置的上下文键名
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		// 返回 401 错误
		return echo.ErrUnauthorized.Wrap(err)
	}

	// 断言 token.Claims 为 *pkg.Claims
	claims, ok := token.Claims.(*pkg.Claims)
	if !ok {
		return echo.ErrUnauthorized.Wrap(errors.New("token claims is not *pkg.Claims"))
	}

	userID := claims.UserID

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, errno.ParmError)
	}

	// 将字符串转换为 uint
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errno.ParmError)
	}

	// 只能查自己：path 的 id 必须和 token 里的用户一致
	if uint(idUint) != userID {
		return echo.ErrForbidden
	}

	// 获取用户信息
	resp, err := u.Svc.GetByID(userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.Success(resp))
}
