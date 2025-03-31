package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 通用结构返回

// Result 结果结构体
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := Result{}
	res.Code = int(ApiCode.Success)
	res.Message = ApiCode.GetMessage(ApiCode.Success)
	res.Data = data
	c.JSON(http.StatusOK, res)
}

// Failed 失败返回
func Failed(c *gin.Context, code int, message string) {
	res := Result{}
	res.Code = code
	res.Message = message
	res.Data = gin.H{}
	c.JSON(http.StatusOK, res)
}
