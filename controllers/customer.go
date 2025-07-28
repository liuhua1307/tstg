package controllers

import (
	"strconv"
	"tangsong-esports/database"
	"tangsong-esports/models"
	"tangsong-esports/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetCustomers 获取客户列表
// @Summary 获取客户列表
// @Description 分页获取客户列表，支持搜索和筛选
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param account query string false "账号搜索"
// @Param customer_name query string false "客户名称搜索"
// @Param phone_number query string false "手机号搜索"
// @Param status query string false "状态筛选"
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/customers [get]
func GetCustomers(c *gin.Context) {
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
	query := database.DB.Model(&models.Customer{})

	// 搜索条件
	if account := c.Query("account"); account != "" {
		query = query.Where("account LIKE ?", "%"+account+"%")
	}
	if customerName := c.Query("customer_name"); customerName != "" {
		query = query.Where("customer_name LIKE ?", "%"+customerName+"%")
	}
	if phoneNumber := c.Query("phone_number"); phoneNumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+phoneNumber+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var customers []models.Customer
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("FinancialInfo").Preload("Preferences").
		Offset(offset).Limit(req.PageSize).Find(&customers).Error

	if err != nil {
		utils.Error(c, "查询失败")
		return
	}

	response := models.PageResponse{
		List:     customers,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.Success(c, response)
}

// CreateCustomer 创建客户
// @Summary 创建客户
// @Description 创建新的客户，包括财务信息和偏好设置
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param customer body models.CustomerCreateRequest true "客户信息"
// @Success 200 {object} models.Response{data=models.Customer}
// @Router /api/v1/customers [post]
func CreateCustomer(c *gin.Context) {
	var req models.CustomerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 检查账号是否已存在
	var existingCustomer models.Customer
	err := database.DB.Where("account = ?", req.Account).First(&existingCustomer).Error
	if err == nil {
		utils.Error(c, "账号已存在")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.Error(c, "查询账号失败")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 处理密码
	var passwordHash string
	var passwordStatus string
	if req.Password != "" {
		// 如果提供了密码，则进行哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			utils.Error(c, "密码加密失败")
			return
		}
		passwordHash = string(hashedPassword)
		passwordStatus = "set"
	} else {
		// 如果没有提供密码，则设置为未设置状态
		passwordHash = ""
		passwordStatus = "unset"
	}

	// 创建客户基本信息
	customer := models.Customer{
		Account:         req.Account,
		PasswordHash:    passwordHash,
		CustomerName:    req.CustomerName,
		ContactMethod:   req.ContactMethod,
		PhoneNumber:     req.PhoneNumber,
		MemberBirthday:  req.MemberBirthday, // 直接使用 *time.Time 类型
		RoomCode:        req.RoomCode,
		AdditionalInfo1: req.AdditionalInfo1,
		AdditionalInfo2: req.AdditionalInfo2,
		AdditionalInfo3: req.AdditionalInfo3,
		Notes:           req.Notes,
		Status:          "正常",
		AccountType:     "admin_created",
		PasswordStatus:  passwordStatus,
	}

	if err := tx.Create(&customer).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建客户失败")
		return
	}
	// 创建财务信息
	financialInfo := models.CustomerFinancialInfo{
		CustomerID:        customer.CustomerID,
		InitialRealCharge: req.InitialRealCharge,
		TotalConsumption:  0,
		TotalRealCharge:   req.InitialRealCharge,
		CurrentBalance:    req.InitialRealCharge,
	}
	if err := tx.Create(&financialInfo).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建财务信息失败")
		return
	}

	// 创建偏好设置
	preferences := models.CustomerPreferences{
		CustomerID:             customer.CustomerID,
		ExclusiveDiscountType:  req.ExclusiveDiscountType,
		ExclusiveDiscountRatio: req.ExclusiveDiscountRatio,
		PlatformBoss:           req.PlatformBoss,
		ExclusiveCS:            req.ExclusiveCS,
	}
	if err := tx.Create(&preferences).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建偏好设置失败")
		return
	}

	// 提交事务
	tx.Commit()

	// 重新查询完整信息
	database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, customer.CustomerID)

	utils.SuccessWithMessage(c, "创建客户成功", customer)
}

// UpdateCustomer 更新客户
// @Summary 更新客户
// @Description 更新客户信息
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Param customer body models.CustomerCreateRequest true "客户信息"
// @Success 200 {object} models.Response{data=models.Customer}
// @Router /api/v1/customers/{id} [put]
func UpdateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的客户ID")
		return
	}

	var req models.CustomerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 更新基本信息
	customer.CustomerName = req.CustomerName
	customer.ContactMethod = req.ContactMethod
	customer.PhoneNumber = req.PhoneNumber
	customer.MemberBirthday = req.MemberBirthday // 直接使用 *time.Time 类型
	customer.RoomCode = req.RoomCode
	customer.AdditionalInfo1 = req.AdditionalInfo1
	customer.AdditionalInfo2 = req.AdditionalInfo2
	customer.AdditionalInfo3 = req.AdditionalInfo3
	customer.Notes = req.Notes

	if err := tx.Save(&customer).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "更新客户失败")
		return
	}

	// 更新偏好设置
	var preferences models.CustomerPreferences
	if err := tx.Where("customer_id = ?", customer.CustomerID).First(&preferences).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			preferences = models.CustomerPreferences{
				CustomerID:             customer.CustomerID,
				ExclusiveDiscountType:  req.ExclusiveDiscountType,
				ExclusiveDiscountRatio: req.ExclusiveDiscountRatio,
				PlatformBoss:           req.PlatformBoss,
				ExclusiveCS:            req.ExclusiveCS,
			}
			tx.Create(&preferences)
		} else {
			tx.Rollback()
			utils.Error(c, "查询偏好设置失败")
			return
		}
	} else {
		preferences.ExclusiveDiscountType = req.ExclusiveDiscountType
		preferences.ExclusiveDiscountRatio = req.ExclusiveDiscountRatio
		preferences.PlatformBoss = req.PlatformBoss
		preferences.ExclusiveCS = req.ExclusiveCS
		tx.Save(&preferences)
	}

	// 提交事务
	tx.Commit()

	// 重新查询完整信息
	database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, customer.CustomerID)

	utils.SuccessWithMessage(c, "更新客户成功", customer)
}

// DeleteCustomer 删除客户
// @Summary 删除客户
// @Description 软删除客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} models.Response
// @Router /api/v1/customers/{id} [delete]
func DeleteCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的客户ID")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 软删除
	if err := database.DB.Delete(&customer).Error; err != nil {
		utils.Error(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除客户成功", nil)
}

// GetCustomerByID 根据ID获取客户
// @Summary 获取客户详情
// @Description 根据ID获取客户详细信息
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} models.Response{data=models.Customer}
// @Router /api/v1/customers/{id} [get]
func GetCustomerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的客户ID")
		return
	}

	var customer models.Customer
	if err := database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}
	utils.Success(c, customer)
}

// RechargeCustomer 客户充值
// @Summary 客户充值
// @Description 为客户账户充值
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Param recharge body models.CustomerRechargeRequest true "充值信息"
// @Success 200 {object} models.Response{data=models.CustomerRechargeHistory}
// @Router /api/v1/customers/{id}/recharge [post]
func RechargeCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的客户ID")
		return
	}

	var req models.CustomerRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 验证客户是否存在
	var customer models.Customer
	if err := database.DB.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询客户失败")
		}
		return
	}

	// 获取操作员ID
	operatorID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到操作员信息")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 计算充值总额
	totalRechargeAmount := req.RealChargeAmount + req.GiftAmount

	// 创建充值记录
	rechargeRecord := models.CustomerRechargeHistory{
		CustomerID:          uint(id),
		RealChargeAmount:    req.RealChargeAmount,
		GiftAmount:          req.GiftAmount,
		TotalRechargeAmount: totalRechargeAmount,
		PaymentMethod:       req.PaymentMethod,
		TransactionID:       req.TransactionID,
		Notes:               req.Notes,
		OperatorID:          operatorID.(uint),
		RechargeAt:          time.Now(),
	}

	if err := tx.Create(&rechargeRecord).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建充值记录失败")
		return
	}

	// 更新客户财务信息
	var financialInfo models.CustomerFinancialInfo
	if err := tx.Where("customer_id = ?", uint(id)).First(&financialInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在财务信息，则创建
			financialInfo = models.CustomerFinancialInfo{
				CustomerID:        uint(id),
				InitialRealCharge: 0,
				TotalConsumption:  0,
				TotalRealCharge:   req.RealChargeAmount,
				CurrentBalance:    totalRechargeAmount,
			}
			tx.Create(&financialInfo)
		} else {
			tx.Rollback()
			utils.Error(c, "查询财务信息失败")
			return
		}
	} else {
		// 更新财务信息
		financialInfo.TotalRealCharge += req.RealChargeAmount
		financialInfo.CurrentBalance += totalRechargeAmount
		tx.Save(&financialInfo)
	}

	// 提交事务
	tx.Commit()
	// 重新查询充值记录
	database.DB.Preload("Customer").Preload("Operator").
		First(&rechargeRecord, rechargeRecord.RechargeID)

	utils.SuccessWithMessage(c, "充值成功", rechargeRecord)
}

// GetCustomerRechargeHistory 获取客户充值记录
// @Summary 获取客户充值记录
// @Description 获取指定客户的充值记录
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/customers/{id}/recharge-history [get]
func GetCustomerRechargeHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的客户ID")
		return
	}

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

	// 验证客户是否存在
	var customer models.Customer
	if err := database.DB.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询客户失败")
		}
		return
	}
	// 查询充值记录
	query := database.DB.Model(&models.CustomerRechargeHistory{}).Where("customer_id = ?", uint(id))

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var rechargeHistory []models.CustomerRechargeHistory
	offset := (req.Page - 1) * req.PageSize
	err = query.Preload("Operator").Order("recharge_at DESC").
		Offset(offset).Limit(req.PageSize).Find(&rechargeHistory).Error

	if err != nil {
		utils.Error(c, "查询充值记录失败")
		return
	}

	response := models.PageResponse{
		List:     rechargeHistory,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.Success(c, response)
}

// AdminResetCustomerPassword 管理员重置客户密码
// @Summary 管理员重置客户密码
// @Description 管理员可以直接为任何客户重置密码
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param reset body models.AdminResetCustomerPasswordRequest true "重置信息"
// @Success 200 {object} models.StandardResponse
// @Router /api/v1/customers/reset-password [post]
func AdminResetCustomerPassword(c *gin.Context) {
	var req models.AdminResetCustomerPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, "密码加密失败")
		return
	}

	// 更新密码和状态
	updates := map[string]interface{}{
		"password_hash":   string(hashedPassword),
		"password_status": "set",
	}

	if err := database.DB.Model(&customer).Updates(updates).Error; err != nil {
		utils.Error(c, "密码重置失败")
		return
	}

	utils.Success(c, gin.H{
		"success": true,
		"message": "客户密码重置成功",
	})
}
