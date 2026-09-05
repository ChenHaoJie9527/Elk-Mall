package service

import (
	"fmt"

	errno "github.com/ChenHaoJie9527/Elk-Mall/internal/common/Errno"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/do"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/dto"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
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
}

// 根据ID获取用户
func (s *UserService) GetByID(id uint) (*dto.UserResp, error) {

}
