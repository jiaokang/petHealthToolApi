package model

import (
	"gorm.io/gorm"
	"time"
)

// 定义任务类型的枚举
type TaskType string

const (
	TaskTypeVaccination  TaskType = "vaccination"   // 疫苗接种
	TaskTypeDeworming    TaskType = "deworming"     // 驱虫
	TaskTypeBath         TaskType = "bath"          // 洗澡
	TaskTypeGrooming     TaskType = "grooming"      // 美容
	TaskTypeMedicalCheck TaskType = "medical_check" // 体检
	TaskTypeMedicine     TaskType = "medicine"      // 喂药
	TaskTypeOther        TaskType = "other"         // 其他
)

// 验证任务类型是否有效
func (t TaskType) IsValid() bool {
	switch t {
	case TaskTypeVaccination, TaskTypeDeworming, TaskTypeBath,
		TaskTypeGrooming, TaskTypeMedicalCheck, TaskTypeMedicine, TaskTypeOther:
		return true
	default:
		return false
	}
}

// Scheduleds 任务排期
type Scheduleds struct {
	gorm.Model
	PetId        uint      `gorm:"index;comment:宠物ID"`
	UserId       uint      `gorm:"index;comment:用户ID"`
	RecordId     uint      `gorm:"comment:记录ID"`
	TaskType     TaskType  `gorm:"size:20;comment:任务类型"`
	ExpectDate   time.Time `gorm:"not null;comment:预计日期"`
	ExecuteDate  time.Time `gorm:"default:NULL;comment:执行日期"`
	ExecuteState bool      `gorm:"comment:执行状态"`
	NotifyState  bool      `gorm:"comment:通知状态"`
}

func (Scheduleds) TableName() string {
	return "scheduleds"
}
