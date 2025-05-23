package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"petHealthToolApi/config"
)

// 颜色
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct {
	Prefix string // 增加前缀字段
}

// Format 实现Formatter(entry *logrus.Entry) ([]byte, error)接口
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同的level去展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		//自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出格式，增加前缀
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s %s\n", timestamp, levelColor, entry.Level, t.Prefix, fileVal, funcVal, entry.Message)
	} else {
		//自定义输出格式，增加前缀
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s\n", timestamp, levelColor, entry.Level, t.Prefix, entry.Message)
	}
	return b.Bytes(), nil
}

var log *logrus.Logger

func init() {
	log = NewLog()
}

func NewLog() *logrus.Logger {
	mLog := logrus.New()                                                  //新建一个实例
	mLog.SetOutput(os.Stdout)                                             //设置输出类型
	mLog.SetReportCaller(true)                                            //开启返回函数名和行号
	mLog.SetFormatter(&LogFormatter{Prefix: config.Config.Logger.Prefix}) //设置自己定义的Formatter，并增加前缀
	mLog.SetLevel(logrus.DebugLevel)                                      //设置最低的Level
	return mLog
}
