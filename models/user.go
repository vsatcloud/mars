package models

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username    string `gorm:"type:varchar(128)"` // 用户名
	Password    string `gorm:"type:varchar(128)"` // 密码
	Salt        string `gorm:"type:varchar(128)"` // 盐
	AuthorityID string `gorm:"type:varchar(128)"` // 身份
	RealName    string `gorm:"type:varchar(128)"` // 真实姓名ar
}
