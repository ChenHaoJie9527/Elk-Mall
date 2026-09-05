package service

import (
	"fmt"
	"time"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/do"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/dto"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepo
	jwtSecret []byte
}

func NewUserService(repo *repository.UserRepo, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

// 注册用户
func (s *UserService) Register(req *dto.RegisterReq) (*dto.UserResp, error) {
	// 查询用户是否存在
	u, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// 如果用户存在，返回用户已存在
	if u != nil {
		return nil, errno.UsernameExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &do.UserDO{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return &dto.UserResp{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

// 登录用户
func (s *UserService) Login(req *dto.LoginReq) (*dto.LoginResp, error) {

	// 查询用户是否存在
	u, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errno.UsernameOrPasswordError
	}

	// 比较密码:明文密码和加密密码比较
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errno.UsernameOrPasswordError
	}

	// 生成 token：NewWithClaims 只是组装未签名的 JWT，SignedString 才会用密钥签出字符串
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,                           // 用户ID
		"exp": now.Add(time.Hour * 24).Unix(), // 过期时间:1天
		"iat": now.Unix(),                     // 签发时间
		"nbf": now.Unix(),                     // 生效时间
		"jti": uuid.New().String(),            // 唯一标识
	})

	// 使用 JWT 密钥签出字符串
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	// 返回登录响应
	return &dto.LoginResp{
		Token: tokenString,
		User: dto.UserResp{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
		},
	}, nil
}

// 根据ID获取用户
func (s *UserService) GetByID(id uint) (*dto.UserResp, error) {
	return nil, nil
}
