// 读取配置文件
package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Config *config

// 总配置文件
type config struct {
	System system `yaml:"system"`
	Logger logger `yaml:"logger"`
	Mysql  mysql  `yaml:"mysql"`
	Redis  redis  `yaml:"redis"`
	Token  token  `yaml:"token"`
	Email  email  `yaml:"mail"`
	Oss    oss    `yaml:"oss"`
}

// 系统配置
type system struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

// 日志
type logger struct {
	Level        string `yaml:"level"`
	Prefix       string `yaml:"prefix"`
	Director     string `yaml:"director"`
	ShowLine     bool   `yaml:"show_line"`
	LogInConsole bool   `yaml:"log_in_console"`
}

// 数据库配置
type mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	LogLevel string `yaml:"log_level"`
	Charset  string `yaml:"charset"`
	MaxIdle  int    `yaml:"maxIdle"`
	MaxOpen  int    `yaml:"maxOpen"`
}

// redis配置
type redis struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

// token配置
type token struct {
	Headers    string `yaml:"headers"`
	ExpireTime int    `yaml:"expire_time"`
	Secret     string `yaml:"secret"`
	Issuer     string `yaml:"issuer"`
}

// 邮箱配置
type email struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

// oss配置
type oss struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Bucket          string `yaml:"bucket"`
	Region          string `yaml:"region"`
}

// 初始化配置
func init() {
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return
	}
	yaml.Unmarshal(yamlFile, &Config)
}
