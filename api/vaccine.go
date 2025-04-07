package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"petHealthToolApi/core"
	"petHealthToolApi/global"
	"petHealthToolApi/model"
	"petHealthToolApi/result"
)

// 疫苗相关

// AddVaccine 新增疫苗记录
func AddVaccine(c *gin.Context) {
	// 1. 获取并验证用户ID
	userId, err := getAndValidateUserId(c)
	if err != nil {
		return
	}

	// 2. 绑定并验证参数
	var param model.AddVaccinationRecord
	if err := bindAndValidate(c, &param); err != nil {
		return
	}

	// 3. 验证宠物权限
	if !hasPetPermission(c, param.PetId, userId) {
		return
	}

	// 4. 在事务中执行操作
	db := core.GetDb()
	err = db.Transaction(func(tx *gorm.DB) error {
		// 创建疫苗记录
		records := param.ToVaccinationRecords()
		if err := tx.Create(records).Error; err != nil {
			global.Log.Errorf("创建疫苗记录失败 petId:%d 用户:%d 错误:%v",
				param.PetId, userId, err)
			return fmt.Errorf("创建疫苗记录失败")
		}

		// 创建定时任务
		if err := createScheduledInTx(tx, records, userId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		handleError(c, int(result.ApiCode.InternalError), err.Error())
		return
	}

	result.Success(c, nil)
}

// 辅助函数分解
func getAndValidateUserId(c *gin.Context) (uint, error) {
	userId, exists := c.Get("userId")
	if !exists {
		msg := "无法获取用户ID"
		global.Log.Error(msg)
		handleError(c, int(result.ApiCode.Unauthorized), msg)
		return 0, fmt.Errorf(msg)
	}

	id, ok := userId.(uint)
	if !ok {
		msg := "用户ID类型无效"
		global.Log.Error(msg)
		handleError(c, int(result.ApiCode.Unauthorized), msg)
		return 0, fmt.Errorf(msg)
	}

	return id, nil
}

func bindAndValidate(c *gin.Context, param *model.AddVaccinationRecord) error {
	if err := c.ShouldBindJSON(param); err != nil {
		global.Log.Errorf("参数绑定失败: %v", err)
		handleError(c, int(result.ApiCode.BadRequest), "无效的请求参数")
		return err
	}

	// 可以在这里添加更多验证逻辑
	return nil
}

func hasPetPermission(c *gin.Context, petId uint, userId uint) bool {
	petIds := getUserPetIds(userId)
	if !hasOptPermission(petId, petIds) {
		global.Log.Errorf("用户无权限操作宠物 petId:%d 用户:%d", petId, userId)
		handleError(c, int(result.ApiCode.Forbidden), "无权限操作该宠物")
		return false
	}
	return true
}

func createScheduledInTx(tx *gorm.DB, r *model.VaccinationRecords, userId uint) error {
	if !r.Notify {
		return nil
	}

	futureDate := r.RecordDate.AddDate(0, 3, 0)
	scheduled := model.Scheduleds{
		PetId:        r.PetId,
		UserId:       userId,
		RecordId:     r.ID,
		TaskType:     model.TaskTypeVaccination, // 修正为疫苗类型
		ExpectDate:   futureDate,
		ExecuteState: false,
		NotifyState:  false,
	}

	if err := tx.Create(&scheduled).Error; err != nil {
		global.Log.Errorf("创建定时任务失败 petId:%d 用户:%d 错误:%v",
			r.PetId, userId, err)
		return fmt.Errorf("创建定时任务失败")
	}
	return nil
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
