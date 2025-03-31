package model

import (
	"gorm.io/gorm"
	"time"
)

// DewormingRecords 宠物驱虫记录
type DewormingRecords struct {
	gorm.Model
	PetId       uint      `gorm:"index;comment:宠物ID"`
	RecordDate  time.Time `gorm:"not null;comment:记录日期"`
	Weight      float64   `gorm:"not null;comment:体重"`
	Medicine    string    `gorm:"size:100;comment:药物"`
	Temperature float64   `gorm:"not null;comment:体温"`
	Age         int       `gorm:"not null;comment:年龄"`
	HealthState string    `gorm:"size:100;comment:健康状态"`
	Remark      string    `gorm:"size:100;comment:备注"`
}

func (DewormingRecords) TableName() string {
	return "deworming_records"
}
