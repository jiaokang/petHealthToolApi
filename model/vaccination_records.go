package model

import (
	"gorm.io/gorm"
	"time"
)

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
}

// AddVaccinationRecord 添加的接种记录
type AddVaccinationRecord struct {
	PetId       uint      `json:"petId"`
	RecordDate  time.Time `json:"recordDate" time_format:"2006-01-02"`
	Weight      float64   `json:"weight"`
	Medicine    string    `json:"medicine"`
	Temperature float64   `json:"temperature"`
	Age         int       `json:"age"`
	HealthState string    `json:"healthState"`
	Remark      string    `json:"remark"`
}

// DeleteVaccinationRecord 添加的接种记录
type DeleteVaccinationRecord struct {
	PetId uint `json:"petId"`
}

func (a *AddVaccinationRecord) ToVaccinationRecords() *VaccinationRecords {
	return &VaccinationRecords{
		PetId:       a.PetId,
		RecordDate:  a.RecordDate,
		Weight:      a.Weight,
		Medicine:    a.Medicine,
		Temperature: a.Temperature,
		Age:         a.Age,
		HealthState: a.HealthState,
		Remark:      a.Remark,
	}
}

func (VaccinationRecords) TableName() string {
	return "vaccination_records"
}
