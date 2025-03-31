package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"petHealthToolApi/core"
	"petHealthToolApi/global"
	"petHealthToolApi/model"
	"petHealthToolApi/result"
	"petHealthToolApi/utils"
	"strings"
	"time"
)

var (
	log = core.NewLog()
)

// 登录相关

// LoginByPwd 通过邮箱和密码登录
func LoginByPwd(c *gin.Context) {
	var param model.LoginByPass
	if err := c.ShouldBindJSON(&param); err != nil {
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	user, err := findUserByEmail(param.Email)
	if err != nil || user.ID == 0 {
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	if !verifyPassword(user.Pwd, param.Password) {
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	token, err := core.GenerateToken(*user)
	if err != nil {
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	result.Success(c, gin.H{
		"nickName": user.NickName,
		"token":    token,
	})
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(c *gin.Context) {
	var param model.SendEmailCode
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Error("Failed to bind JSON: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	// 生成六位随机数字验证码
	verifyCode := utils.GenerateSixDigitCode()

	err := utils.SendVerifyCode(param.Email, verifyCode)
	if err != nil {
		log.Error("Failed to send email: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	redisDb := core.GetRedisDb()
	if err := redisDb.Set(global.Ctx, param.Email, verifyCode, 5*60*time.Second).Err(); err != nil {
		log.Error("Failed to set email code in Redis: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	result.Success(c, gin.H{
		"message": "验证码发送成功",
	})
}

// LoginByCode 通过邮箱和验证码登录
func LoginByCode(c *gin.Context) {
	var param model.LoginByCode
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Error("Failed to bind JSON: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	redisDb := core.GetRedisDb()
	exists, _ := redisDb.Exists(global.Ctx, param.Email).Result()
	// 根据结果输出
	if exists != 1 {
		log.Error(" email code is expire")
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	if code, err := redisDb.Get(global.Ctx, param.Email).Result(); err != nil {
		log.Error("Failed to get email code from Redis: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	} else if code != param.VerifyCode {
		log.Error("Invalid email code: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	user, err := findUserByEmail(param.Email)
	if err != nil || user.ID == 0 {
		// 注册新用户
		user, err = registerNewUser(param.Email)
	}
	token, err := core.GenerateToken(*user)
	if err != nil {
		log.Error("Failed to generate token: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
	result.Success(c, gin.H{
		"nickName": user.NickName,
		"token":    token,
	})
	if err := redisDb.Del(global.Ctx, param.Email).Err(); err != nil {
		log.Error("Failed to delete email code from Redis: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(501))
		return
	}
}

func findUserByEmail(email string) (*model.Users, error) {
	user := &model.Users{}
	db := core.GetDb()
	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 注册新用户
func registerNewUser(email string) (*model.Users, error) {
	// 提取邮箱前缀作为昵称
	nickName := strings.Split(email, "@")[0]

	// 创建用户对象
	user := &model.Users{
		NickName: nickName,
		Email:    email,
	}

	// 插入数据库
	db := core.GetDb()
	if err := db.Create(user).Error; err != nil {
		log.Error("Failed to create user: %v", err)
		return nil, err
	}

	// 返回注册成功的用户
	return user, nil
}

func verifyPassword(hashedPassword, inputPassword string) bool {
	// 这里可以使用 bcrypt 或其他加密算法来验证密码
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)) == nil
}

func handleError(c *gin.Context, code int, message string) {
	result.Failed(c, code, message)
}
