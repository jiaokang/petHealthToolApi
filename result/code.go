package result

// Codes 定义状态
type Codes struct {
	Message map[uint]string
	Success uint
	Failed  uint
	// HTTP 标准状态码
	BadRequest       uint
	Unauthorized     uint
	Forbidden        uint
	NotFound         uint
	MethodNotAllowed uint
	// 业务状态码
	InvalidParams   uint
	RecordExists    uint
	RecordNotExists uint
	// 系统状态码
	InternalError uint
	DatabaseError uint
	CacheError    uint
	// 限流相关
	TooManyRequests uint
}

// ApiCode 定义状态码和状态信息
var ApiCode = &Codes{
	Success:          200,
	Failed:           501,
	BadRequest:       400,
	Unauthorized:     401,
	Forbidden:        403,
	NotFound:         404,
	MethodNotAllowed: 405,
	InvalidParams:    10001,
	RecordExists:     10002,
	RecordNotExists:  10003,
	InternalError:    500,
	DatabaseError:    5001,
	CacheError:       5002,
	TooManyRequests:  429,
}

// 状态信息初始化
func init() {
	ApiCode.Message = map[uint]string{
		ApiCode.Success:          "成功",
		ApiCode.Failed:           "失败",
		ApiCode.BadRequest:       "请求参数错误",
		ApiCode.Unauthorized:     "未授权，请登录",
		ApiCode.Forbidden:        "禁止访问",
		ApiCode.NotFound:         "资源不存在",
		ApiCode.MethodNotAllowed: "方法不允许",
		ApiCode.InvalidParams:    "参数校验失败",
		ApiCode.RecordExists:     "记录已存在",
		ApiCode.RecordNotExists:  "记录不存在",
		ApiCode.InternalError:    "服务器内部错误",
		ApiCode.DatabaseError:    "数据库操作失败",
		ApiCode.CacheError:       "缓存操作失败",
		ApiCode.TooManyRequests:  "请求过于频繁，请稍后再试",
	}
}

// GetMessage 获取状态信息
func (c *Codes) GetMessage(code uint) string {
	message, ok := c.Message[code]
	if !ok {
		return "未知状态"
	}
	return message
}
