package core

import (
	"github.com/golang-jwt/jwt/v5"
	"petHealthToolApi/config"
	"petHealthToolApi/model"
	"time"
)

var (
	TokenExpireDuration = time.Duration(config.Config.Token.ExpireTime) * time.Hour // token 过期时间
	Secret              = []byte(config.Config.Token.Secret)                        // 密钥
	Issuer              = config.Config.Token.Issuer                                // 签发人
)

// GenerateToken 生成 JWT Token
func GenerateToken(user model.Users) (string, error) {
	claims := model.UserStdClaim{
		JwtUser: model.JwtUser{
			Id:       user.ID,
			NickName: user.NickName,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			Issuer:    Issuer,                                                  // 签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Secret)
}
