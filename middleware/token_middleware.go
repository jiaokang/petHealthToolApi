package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"petHealthToolApi/config"
	"petHealthToolApi/model"
	"strings"
)

// AuthMiddleware 验证 Token 并提取用户信息
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中提取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// 检查 Token 格式
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// 解析和验证 Token
		token, err := jwt.ParseWithClaims(tokenString, &model.UserStdClaim{}, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			// 返回签名密钥
			return []byte(config.Config.Token.Secret), nil
		})
		if err != nil {
			log.Printf("Failed to parse token: %v", err) // 打印错误信息
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 提取用户信息
		if claims, ok := token.Claims.(*model.UserStdClaim); ok && token.Valid {
			c.Set("userId", claims.Id)         // 将 userID 存入上下文
			c.Set("nickName", claims.NickName) // 将 nickName 存入上下文
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Next()
	}
}
