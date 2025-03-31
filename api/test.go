package api

import (
	"github.com/gin-gonic/gin"
	"petHealthToolApi/result"
)

// 返回测试

// Success 测试成功
func Success(c *gin.Context) {
	result.Success(c, 200)
}

// Failed 测试失败
func Failed(c *gin.Context) {
	result.Failed(c, 500, "测试失败")
}
