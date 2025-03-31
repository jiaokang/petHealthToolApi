package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"petHealthToolApi/core"
	"petHealthToolApi/global"
	"petHealthToolApi/model"
	"petHealthToolApi/result"
)

// 个人中心

// SetPassword 设置密码
func SetPassword(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	var setPwdParam model.SetPasswordReq
	if err := c.ShouldBindJSON(&setPwdParam); err != nil {
		global.Log.Error("Failed to bind JSON: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	pwdHash, err := hashPassword(setPwdParam.Password)
	if err != nil {
		global.Log.Error("Failed to hash password: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	db := core.GetDb()
	if err := db.Model(&model.Users{}).Where("id = ?", userId).Update("pwd", pwdHash).Error; err != nil {
		global.Log.Error("Failed to update password: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	result.Success(c, nil)
}

// hashPassword 使用 bcrypt 算法对密码进行加密
// cost 参数表示加密的计算成本（4-31），默认推荐使用 bcrypt.DefaultCost(10)
func hashPassword(password string, cost ...int) (string, error) {
	// 设置默认的 cost 值
	bcryptCost := bcrypt.DefaultCost
	if len(cost) > 0 {
		bcryptCost = cost[0]
	}
	// 生成密码的哈希值
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
