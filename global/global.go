package global

import (
	"context"
	"github.com/sirupsen/logrus"
)

// 全局共享配置

var (
	Log *logrus.Logger
	Ctx = context.Background()
)
