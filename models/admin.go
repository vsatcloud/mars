package models

import (
	"github.com/vsatcloud/mars/utils"
	"gorm.io/gorm"
)

type AdminModel struct {
	gorm.Model
	Username string `gorm:"type:varchar(128)"` // 用户名
	Password string `gorm:"type:varchar(128)"` // 密码
	Salt     string `gorm:"type:varchar(128)"` // 盐
	RealName string `gorm:"type:varchar(128)"` // 真实姓名ar
	Roles    string `gorm:"type:varchar(512)"` // 角色
	Avatar   string //头像
}

func (am *AdminModel) GenToken(secret string) (string, error) {
	var data = map[string]interface{}{
		"user_id":    am.ID,
		"username":   am.Username,
		"roles":      am.Roles,
		"created_at": am.CreatedAt.Unix(),
	}

	return utils.GenerateToken(data, secret)
}
