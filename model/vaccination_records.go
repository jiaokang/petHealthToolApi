package model

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type Date time.Time

// VaccinationRecords 接种记录
type VaccinationRecords struct {
	gorm.Model
	PetId       uint      `gorm:"index;comment:宠物ID"`
	RecordDate  time.Time `gorm:"not null;comment:记录日期"`
	Weight      float64   `gorm:"not null;comment:体重"`
	Medicine    string    `gorm:"size:100;comment:药物"`
	Temperature float64   `gorm:"not null;comment:体温"`
	Age         int       `gorm:"not null;comment:年龄"`
	HealthState string    `gorm:"size:100;comment:健康状态"`
	Remark      string    `gorm:"size:100;comment:备注"`
	Notify      bool      `gorm:"comment:通知状态"`
}

// AddVaccinationRecord 添加的接种记录
type AddVaccinationRecord struct {
	PetId       uint    `json:"petId" validate:"required"`
	RecordDate  Date    `json:"recordDate" validate:"required"`
	Weight      float64 `json:"weight" validate:"required"`
	Medicine    string  `json:"medicine" validate:"required"`
	Temperature float64 `json:"temperature" validate:"required"`
	Age         int     `json:"age" validate:"required"`
	HealthState string  `json:"healthState" validate:"required"`
	Remark      string  `json:"remark"`
	Notify      bool    `json:"notify" validate:"required"`
}

// DeleteVaccinationRecord 添加的接种记录
type DeleteVaccinationRecord struct {
	PetId uint `json:"petId"`
}

func (a *AddVaccinationRecord) ToVaccinationRecords() *VaccinationRecords {
	// 将 Date 类型转换为 time.Time
	recordDate := time.Time(a.RecordDate)
	return &VaccinationRecords{
		PetId:       a.PetId,
		RecordDate:  recordDate,
		Weight:      a.Weight,
		Medicine:    a.Medicine,
		Temperature: a.Temperature,
		Age:         a.Age,
		HealthState: a.HealthState,
		Remark:      a.Remark,
		Notify:      a.Notify,
	}
}

func (VaccinationRecords) TableName() string {
	return "vaccination_records"
}

// 自定义 JSON 解析
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s) // 指定你的格式
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}
