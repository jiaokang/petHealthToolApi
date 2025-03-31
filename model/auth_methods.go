package model

import "gorm.io/gorm"

// AuthMethods 认证方式
type AuthMethods struct {
	gorm.Model
	UserId     uint   `gorm:"index;comment:用户ID"`
	AuthType   string `gorm:"size:20;comment:认证类型"`
	AuthValue  string `gorm:"size:100;comment:认证值"`
	Credential string `gorm:"size:100;comment:认证凭据"`
}

func (AuthMethods) TableName() string {
	return "auth_methods"
}
