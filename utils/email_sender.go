package utils

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"math/rand"
	"os"
	"petHealthToolApi/config"
	"time"
)

// SendVerifyCode 发送验证码
func SendVerifyCode(emailTo string, verifyCode string) error {
	cfg := config.Config.Email

	// 配置邮件发送者信息
	from := cfg.User
	password := cfg.Pass
	to := emailTo
	smtpHost := cfg.Host
	smtpPort := cfg.Port

	// 读取 HTML 文件
	htmlContent, err := os.ReadFile("./verifycode_email_template.html")
	if err != nil {
		log.Fatal("Failed to read HTML file:", err)
		return err
	}

	// 定义替换的数据
	data := struct {
		VerificationCode string
		ExpirationTime   string
		CurrentYear      string
	}{
		VerificationCode: verifyCode,                // 替换为实际的验证码
		ExpirationTime:   "5",                       // 替换为实际的过期时间
		CurrentYear:      time.Now().Format("2006"), // 获取当前年份
	}

	// 解析 HTML 模板并替换占位符
	tmpl, err := template.New("email").Parse(string(htmlContent))
	if err != nil {
		log.Fatal("Failed to parse HTML template:", err)
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Fatal("Failed to execute template:", err)
		return err
	}

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "验证码邮件")
	m.SetBody("text/html", body.String())

	// 创建 SMTP 客户端
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Failed to send email:", err)
		return err
	}
	log.Println("Email sent successfully!")
	return nil
}

// GenerateSixDigitCode 生成一个六位数字验证码
func GenerateSixDigitCode() string {
	// 设置随机种子，确保每次运行生成的验证码不同
	rand.Seed(time.Now().UnixNano())

	// 生成一个 0 到 999999 之间的随机数
	code := rand.Intn(1000000)

	// 将随机数格式化为六位数字字符串，不足六位时前面补零
	return fmt.Sprintf("%06d", code)
}
