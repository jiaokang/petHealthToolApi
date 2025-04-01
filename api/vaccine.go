package api

import (
	"github.com/gin-gonic/gin"
	"petHealthToolApi/core"
	"petHealthToolApi/global"
	"petHealthToolApi/model"
	"petHealthToolApi/result"
)

// 疫苗相关

// AddVaccine 新增疫苗记录
func AddVaccine(c *gin.Context) {
	// 1. 获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Unauthorized), result.ApiCode.GetMessage(result.ApiCode.Unauthorized))
		return
	}

	var addVaccineParam model.AddVaccinationRecord
	if err := c.ShouldBindJSON(&addVaccineParam); err != nil {
		global.Log.Error("Failed to bind JSON: %v", err)
		handleError(c, int(result.ApiCode.InternalError), result.ApiCode.GetMessage(result.ApiCode.InternalError))
		return
	}
	petIds := getUserPetIds(userId)
	if !hasOptPermission(addVaccineParam.PetId, petIds) {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Forbidden), result.ApiCode.GetMessage(result.ApiCode.Forbidden))
		return
	}

	vaccinationRecords := addVaccineParam.ToVaccinationRecords()
	db := core.GetDb()
	if err := db.Create(&vaccinationRecords); err != nil {
		global.Log.Error("Failed to create vaccination record: %v", err)
		handleError(c, int(result.ApiCode.InternalError), result.ApiCode.GetMessage(result.ApiCode.InternalError))
		return
	}
	result.Success(c, nil)
}

// GetVaccineList 获取疫苗记录
func GetVaccineList(c *gin.Context) {
	// 1. 获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Unauthorized), result.ApiCode.GetMessage(result.ApiCode.Unauthorized))
		return
	}
	// 2. 获取查询参数
	petId := c.Query("petId")
	db := core.GetDb()
	// 3. 准备查询条件
	query := db.Where("user_id = ?", userId)
	if petId != "" {
		query = query.Where("id = ?", petId)
	}
	// 4. 查询宠物列表
	var pets []model.Pets
	if err := query.Find(&pets).Error; err != nil {
		global.Log.Error("Failed to get pet list: %v", err)
		handleError(c, int(result.ApiCode.InternalError), "获取宠物列表失败")
		return
	}
	// 5. 如果没有宠物直接返回空数组
	if len(pets) == 0 {
		result.Success(c, []model.VaccinationRecords{})
		return
	}
	// 6. 收集所有宠物ID
	petIds := make([]uint, 0, len(pets))
	for _, pet := range pets {
		petIds = append(petIds, pet.ID)
	}
	// 7. 批量查询疫苗记录
	var vaccinationRecords []model.VaccinationRecords
	if err := db.Where("pet_id IN (?)", petIds).Find(&vaccinationRecords).Error; err != nil {
		global.Log.Error("Failed to get vaccination records: %v", err)
		handleError(c, int(result.ApiCode.InternalError), "获取疫苗记录失败")
		return
	}
	// 8. 返回结果
	result.Success(c, vaccinationRecords)
}

// DeleteVaccine 删除疫苗记录
func DeleteVaccine(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Unauthorized), result.ApiCode.GetMessage(result.ApiCode.Unauthorized))
		return
	}

	var deleteVaccineParam model.DeleteVaccinationRecord
	if err := c.ShouldBindJSON(&deleteVaccineParam); err != nil {
		global.Log.Error("Failed to bind JSON: %v", err)
		handleError(c, int(result.ApiCode.InternalError), result.ApiCode.GetMessage(result.ApiCode.InternalError))
		return
	}

	petIds := getUserPetIds(userId)
	if !hasOptPermission(deleteVaccineParam.PetId, petIds) {
		global.Log.Error("Failed to get userId from context")
		handleError(c, int(result.ApiCode.Forbidden), result.ApiCode.GetMessage(result.ApiCode.Forbidden))
		return
	}
	db := core.GetDb()
	if err := db.Where("pet_id = ?", deleteVaccineParam.PetId).Delete(&model.VaccinationRecords{}).Error; err != nil {
		global.Log.Error("Failed to delete vaccination record: %v", err)
		handleError(c, int(result.ApiCode.InternalError), result.ApiCode.GetMessage(result.ApiCode.InternalError))
		return
	}
	result.Success(c, nil)
}

// hasOptPermission 检查操作的宠物ID是否在用户宠物列表中
func hasOptPermission(optPetId uint, petIds []uint) bool {
	for _, petId := range petIds {
		if optPetId == petId {
			return true
		}
	}
	return false
}

// getUserPetIds 获取用户宠物列表
func getUserPetIds(userId any) []uint {
	db := core.GetDb()
	var pets []model.Pets
	if err := db.Where("user_id = ?", userId).Find(&pets).Error; err != nil {
		global.Log.Error("Failed to get pet list: %v", err)
		return nil
	}
	petIds := make([]uint, 0, len(pets))
	for _, pet := range pets {
		petIds = append(petIds, pet.ID)
	}
	return petIds
}
