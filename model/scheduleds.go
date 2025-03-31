package model

import (
	"gorm.io/gorm"
	"time"
)

// Scheduleds 任务排期
type Scheduleds struct {
	gorm.Model
	PetId        uint      `gorm:"index;comment:宠物ID"`
	UserId       uint      `gorm:"index;comment:用户ID"`
	TaskType     string    `gorm:"size:20;comment:任务类型"`
	ExpectDate   time.Time `gorm:"not null;comment:预计日期"`
	ExecuteDate  time.Time `gorm:"comment:执行日期"`
	ExecuteState bool      `gorm:"comment:执行状态"`
	NotiftyState bool      `gorm:"comment:通知状态"`
}

func (Scheduleds) TableName() string {
	return "scheduleds"
}
