package router

import (
	"github.com/gin-gonic/gin"
	"petHealthToolApi/api"
	"petHealthToolApi/config"
	"petHealthToolApi/middleware"
)

// 初始化路由以及注册路由

func InitRouter() *gin.Engine {
	// 设置启动模式
	gin.SetMode(config.Config.System.Env)
	router := gin.New()
	// 宕机恢复
	router.Use(gin.Recovery())
	// 注册路由
	register(router)
	return router
}

// register 路由注册接口
func register(router *gin.Engine) {
	// 认证相关
	authGroup := router.Group("/api/auth")
	{
		//  使用邮箱➕密码登录
		authGroup.POST("/loginByPwd", api.LoginByPwd)
		// 发送邮箱验证码
		authGroup.POST("/sendEmailCode", api.SendEmailCode)
		// 使用邮箱+验证码登录
		authGroup.POST("/loginByCode", api.LoginByCode)
	}

	// 个人中心
	profileGroup := router.Group("/api/profile")
	{
		profileGroup.Use(middleware.AuthMiddleware())
		//
		profileGroup.PUT("/setPwd", api.SetPassword)
	}

	// 宠物相关
	petGroup := router.Group("/api/pet")
	{
		petGroup.Use(middleware.AuthMiddleware())
		petGroup.POST("/addPet", api.AddPet)
		petGroup.GET("/getPetList", api.GetPetList)
		petGroup.DELETE("/deletePet", api.DeletePet)

	}

}
