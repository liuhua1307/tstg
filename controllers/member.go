package controllers

import (
	"log"
	"strconv"
	"tangsong-esports/database"
	"tangsong-esports/models"
	"tangsong-esports/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetMembers 获取成员列表
// @Summary 获取内部成员列表
// @Description 分页获取内部成员列表，支持搜索和筛选
// @Tags 成员管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param account query string false "账号搜索"
// @Param name query string false "姓名搜索"
// @Param department query string false "部门筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/members [get]
func GetMembers(c *gin.Context) {
	var req models.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询
	query := database.DB.Model(&models.InternalMember{})

	// 搜索条件
	if account := c.Query("account"); account != "" {
		query = query.Where("account LIKE ?", "%"+account+"%")
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if department := c.Query("department"); department != "" {
		query = query.Where("department = ?", department)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var members []models.InternalMember
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Permissions").Preload("FinancialSettings").Preload("Relationships").
		Offset(offset).Limit(req.PageSize).Find(&members).Error

	if err != nil {
		utils.Error(c, "查询失败")
		return
	}

	response := models.PageResponse{
		List:     members,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.Success(c, response)
}

// CreateMember 创建成员
// @Summary 创建内部成员
// @Description 创建新的内部成员，包括权限和财务设置
// @Tags 成员管理
// @Accept json
// @Produce json
// @Param member body models.MemberCreateRequest true "成员信息"
// @Success 200 {object} models.Response{data=models.InternalMember}
// @Router /api/v1/members [post]
func CreateMember(c *gin.Context) {
	var req models.MemberCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 检查账号是否已存在
	var existingMember models.InternalMember
	err := database.DB.Where("account = ?", req.Account).First(&existingMember).Error
	if err == nil {
		utils.Error(c, "账号已存在")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.Error(c, "查询账号失败")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(c, "密码加密失败")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 创建成员基本信息
	member := models.InternalMember{
		Account:      req.Account,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		Department:   req.Department,
		UserRole:     req.UserRole,
		Status:       "正常",
		IsEnabled:    true,
		Notes:        req.Notes,
	}

	if err := tx.Create(&member).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建成员失败")
		return
	}

	// 创建权限设置
	permissions := models.MemberPermissions{
		MemberID:       member.MemberID,
		IsAuditor:      req.IsAuditor,
		CanReport:      req.CanReport,
		CanAcceptOrder: req.CanAcceptOrder,
	}
	if err := tx.Create(&permissions).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建权限设置失败")
		return
	}

	// 创建财务设置
	financialSettings := models.MemberFinancialSettings{
		MemberID:       member.MemberID,
		CommissionRate: req.CommissionRate,
	}
	if err := tx.Create(&financialSettings).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建财务设置失败")
		return
	}

	// 创建关系设置
	relationships := models.MemberRelationships{
		MemberID:   member.MemberID,
		CreatorID:  req.CreatorID,
		AssigneeID: req.AssigneeID,
	}
	if req.AssigneeID != nil {
		now := time.Now()
		relationships.AssignedAt = &now
	}
	if err := tx.Create(&relationships).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建关系设置失败")
		return
	}

	// 提交事务
	tx.Commit()

	// 重新查询完整信息
	database.DB.Preload("Permissions").Preload("FinancialSettings").Preload("Relationships").
		First(&member, member.MemberID)

	utils.SuccessWithMessage(c, "创建成员成功", member)
}

// UpdateMember 更新成员
// @Summary 更新内部成员
// @Description 更新内部成员信息
// @Tags 成员管理
// @Accept json
// @Produce json
// @Param id path int true "成员ID"
// @Param member body models.MemberCreateRequest true "成员信息"
// @Success 200 {object} models.Response{data=models.InternalMember}
// @Router /api/v1/members/{id} [put]
func UpdateMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的成员ID")
		return
	}

	var req models.MemberCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找成员
	var member models.InternalMember
	if err := database.DB.First(&member, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "成员不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 更新基本信息
	member.Name = req.Name
	member.PhoneNumber = req.PhoneNumber
	member.Department = req.Department
	member.UserRole = req.UserRole
	member.Notes = req.Notes

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			tx.Rollback()
			utils.Error(c, "密码加密失败")
			return
		}
		member.PasswordHash = hashedPassword
	}

	if err := tx.Save(&member).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "更新成员失败")
		return
	}

	// 更新权限设置
	var permissions models.MemberPermissions
	if err := tx.Where("member_id = ?", member.MemberID).First(&permissions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在则创建
			permissions = models.MemberPermissions{
				MemberID:       member.MemberID,
				IsAuditor:      req.IsAuditor,
				CanReport:      req.CanReport,
				CanAcceptOrder: req.CanAcceptOrder,
			}
			tx.Create(&permissions)
		} else {
			tx.Rollback()
			utils.Error(c, "查询权限设置失败")
			return
		}
	} else {
		permissions.IsAuditor = req.IsAuditor
		permissions.CanReport = req.CanReport
		permissions.CanAcceptOrder = req.CanAcceptOrder
		tx.Save(&permissions)
	}

	// 更新财务设置
	var financialSettings models.MemberFinancialSettings
	if err := tx.Where("member_id = ?", member.MemberID).First(&financialSettings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			financialSettings = models.MemberFinancialSettings{
				MemberID:       member.MemberID,
				CommissionRate: req.CommissionRate,
			}
			tx.Create(&financialSettings)
		} else {
			tx.Rollback()
			utils.Error(c, "查询财务设置失败")
			return
		}
	} else {
		financialSettings.CommissionRate = req.CommissionRate
		tx.Save(&financialSettings)
	}

	// 提交事务
	tx.Commit()

	// 重新查询完整信息
	database.DB.Preload("Permissions").Preload("FinancialSettings").Preload("Relationships").
		First(&member, member.MemberID)

	utils.SuccessWithMessage(c, "更新成员成功", member)
}

// DeleteMember 删除成员
// @Summary 删除内部成员
// @Description 软删除内部成员
// @Tags 成员管理
// @Accept json
// @Produce json
// @Param id path int true "成员ID"
// @Success 200 {object} models.Response
// @Router /api/v1/members/{id} [delete]
func DeleteMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的成员ID")
		return
	}

	// 查找成员
	var member models.InternalMember
	if err := database.DB.First(&member, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "成员不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 软删除
	if err := database.DB.Delete(&member).Error; err != nil {
		utils.Error(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成员成功", nil)
}

// GetMemberByID 根据ID获取成员
// @Summary 获取成员详情
// @Description 根据ID获取内部成员详细信息
// @Tags 成员管理
// @Accept json
// @Produce json
// @Param id path int true "成员ID"
// @Success 200 {object} models.Response{data=models.InternalMember}
// @Router /api/v1/members/{id} [get]
func GetMemberByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的成员ID")
		return
	}

	var member models.InternalMember
	if err := database.DB.Preload("Permissions").Preload("FinancialSettings").
		Preload("Relationships").Preload("Relationships.Creator").
		Preload("Relationships.Assignee").First(&member, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "成员不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	utils.Success(c, member)
}
