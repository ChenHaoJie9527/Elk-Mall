package do

import (
	"time"

	"gorm.io/gorm"
)

type UserDO struct {
	ID        uint   `gorm:"primaryKey"`          // 用户ID是主键，自动递增，unit 表示无符号整数，64 表示长度为 64
	Username  string `gorm:"uniqueIndex;size:64"` // 用户名是唯一索引，长度为 64
	Password  string `gorm:"size:128"`            // 密码是长度为 128 的字符串
	Nickname  string `gorm:"size:64"`             // 昵称是长度为 64 的字符串
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 删除时间索引
}
