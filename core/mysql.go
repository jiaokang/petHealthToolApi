package core

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"petHealthToolApi/config"
	"petHealthToolApi/global"
	"petHealthToolApi/model"
)

// mysql初始化配置

var Db *gorm.DB

func init() {
	global.Log = NewLog()
	var err error
	var dbConfig = config.Config.Mysql
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Db, dbConfig.Charset)
	Db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return
	}
	if Db.Error != nil {
		return
	}
	sqlDb, err := Db.DB()
	sqlDb.SetMaxIdleConns(dbConfig.MaxIdle)
	sqlDb.SetMaxOpenConns(dbConfig.MaxOpen)

	// 迁移表结构
	// 自动迁移表结构
	err = Db.AutoMigrate(&model.Users{}, &model.AuthMethods{}, &model.Pets{}, &model.VaccinationRecords{}, &model.Scheduleds{}, &model.DewormingRecords{})
	if err != nil {
		return
	}
	global.Log.Infof("mysql init success")
	return
}

// GetDb 获取db
func GetDb() *gorm.DB {
	if Db == nil {
		panic("mysql not init")
	}
	return Db
}
