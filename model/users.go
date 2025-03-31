package model

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Users 用户信息
type Users struct {
	gorm.Model
	NickName string `gorm:"size:64;comment:昵称"`
	Phone    string `gorm:"size:20;uniqueIndex;comment:主人手机号"`
	Email    string `gorm:"size:100;uniqueIndex;comment:主人邮箱"`
	Address  string `gorm:"type:text;comment:主人地址"`
	Pwd      string `gorm:"size:64;comment:密码"`
}

func (Users) TableName() string {
	return "users"
}

// LoginByPass 邮箱+密码
type LoginByPass struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SendEmailCode 发送验证码
type SendEmailCode struct {
	Email string `json:"email"`
}

// LoginByCode 邮箱+验证码
type LoginByCode struct {
	Email      string `json:"email"`
	VerifyCode string `json:"verifyCode"`
}

// SetPasswordReq 设置密码
type SetPasswordReq struct {
	Password string `json:"password"`
}

// JwtUser jwt用户信息
type JwtUser struct {
	Id         uint   `json:"id"`
	NickName   string `json:"nickName"`
	ExpireTime int64  `json:"expireTime"`
}

// UserStdClaim 自定义的 JWT Claims
type UserStdClaim struct {
	JwtUser
	jwt.RegisteredClaims
}
