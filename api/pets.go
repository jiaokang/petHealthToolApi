package api

import (
	"github.com/gin-gonic/gin"
	"petHealthToolApi/core"
	"petHealthToolApi/global"
	"petHealthToolApi/model"
	"petHealthToolApi/result"
)

// 宠物管理

// AddPet 新增宠物
func AddPet(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	var param model.AddPet
	if err := c.ShouldBindJSON(&param); err != nil {
		global.Log.Error("Failed to bind JSON: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	db := core.GetDb()
	if err := db.Create(&model.Pets{
		UserId:   userId.(uint),
		Name:     param.Name,
		Breed:    param.Breed,
		Birthday: param.Birthday,
		Sex:      param.Sex,
		Avatar:   param.Avatar,
	}).Error; err != nil {
		global.Log.Error("Failed to create pet:%v", err)
	}
	result.Success(c, nil)
}

// GetPetList 获取宠物列表
func GetPetList(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	db := core.GetDb()
	var pets []model.Pets
	if err := db.Where("user_id = ?", userId).Find(&pets).Error; err != nil {
		global.Log.Error("Failed to get pet list: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	// 转换为 PetResponse 结构体，格式化日期
	var response []model.PetResponse
	for _, pet := range pets {
		response = append(response, model.PetResponse{
			PetId:    pet.ID,
			Name:     pet.Name,
			Breed:    pet.Breed,
			Birthday: pet.Birthday.Format("2006-01-02"), // 格式化为 YYYY-MM-DD
			Avatar:   pet.Avatar,
			Sex:      pet.Sex,
		})
	}
	result.Success(c, gin.H{
		"pets": response,
	})
}

// DeletePet 删除宠物
func DeletePet(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	var deletePet model.DeletePet
	if err := c.ShouldBindJSON(&deletePet); err != nil {
		global.Log.Error("Failed to bind JSON: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	db := core.GetDb()
	if err := db.Where("id = ? and user_id=? ", deletePet.PetId, userId).Delete(&model.Pets{}, deletePet.PetId); err != nil {
		global.Log.Error("Failed to delete pet: %v", err)
		handleError(c, int(result.ApiCode.Failed), result.ApiCode.GetMessage(result.ApiCode.Failed))
		return
	}
	result.Success(c, nil)
}
