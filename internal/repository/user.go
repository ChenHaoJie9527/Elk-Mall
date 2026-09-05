package repository

import (
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/do"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

// 创建用户仓库实例
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// 创建用户
func (u *UserRepo) CreateUser(user *do.UserDO) error {
	// 按 UserDO 往对应表 INSERT 一条记录
	// .Error：如果这次操作没有失败则是 nil，失败则是具体错误
	return u.db.Create(user).Error
}

// 根据主键 id 查询用户
func (u *UserRepo) GetByID(id uint) (*do.UserDO, error) {
	var user do.UserDO
	//按照主键 id 查询一条记录
	// .Error：如果这次操作没有失败则是 nil，失败则是具体错误
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	// 返回查询到的记录
	return &user, nil
}

// 获取用户列表
// page: 页码
// pageSize: 每页条数
// 返回: 用户列表, 总条数, 错误
func (u *UserRepo) GetList(page, pageSize int) ([]do.UserDO, int64, error) {
	var total int64
	var list []do.UserDO

	// Model: 指定要查询的表
	// Count: 统计总条数
	u.db.Model(&do.UserDO{}).Count(&total)

	// Offset: 偏移量
	// Limit: 限制条数
	// Find: 查询结果
	if err := u.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// 更新用户
func (u *UserRepo) UpdateUser(user *do.UserDO) error {
	return u.db.Save(user).Error
}

// 删除用户
func (u *UserRepo) DeleteUser(id uint) error {
	return u.db.Delete(&do.UserDO{}, id).Error
}
