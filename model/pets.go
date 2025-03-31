package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// Pets 宠物信息
type Pets struct {
	gorm.Model
	Name     string    `gorm:"size:20;comment:宠物名称"`
	Breed    string    `gorm:"size:50;comment:宠物品种"`
	Birthday time.Time `gorm:"not null;comment:宠物出生日期"`
	UserId   uint      `gorm:"index;comment:用户ID"`
	Avatar   string    `gorm:"size:100;comment:宠物头像"`
	Sex      string    `gorm:"size:10;comment:宠物性别"`
}

func (Pets) TableName() string {
	return "pets"
}

// AddPet 新增宠物
type AddPet struct {
	Name     string    `json:"name"`
	Breed    string    `json:"breed"`
	Birthday time.Time `json:"birthday"`
	Avatar   string    `json:"avatar"`
	Sex      string    `json:"sex"`
}

// PetResponse 宠物响应
type PetResponse struct {
	PetId    uint   `json:"petId"`
	Name     string `json:"name"`
	Breed    string `json:"breed"`
	Birthday string `json:"birthday"` // 格式化为 YYYY-MM-DD
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
}

// DeletePet 删除宠物
type DeletePet struct {
	PetId uint `json:"petId"`
}

// UnmarshalJSON 自定义方法
func (a *AddPet) UnmarshalJSON(data []byte) error {
	type Alias AddPet
	aux := &struct {
		Birthday string `json:"birthday"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// 解析 YYYY-MM-DD 格式的日期
	birthday, err := time.Parse("2006-01-02", aux.Birthday)
	if err != nil {
		return err
	}
	a.Birthday = birthday
	return nil
}
