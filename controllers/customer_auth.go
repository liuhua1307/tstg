package controllers

import (
	"log"
	"tangsong-esports/database"
	"tangsong-esports/models"
	"tangsong-esports/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CustomerRegister 客户注册
// @Summary 客户注册
// @Description 客户自助注册接口
// @Tags 客户认证
// @Accept json
// @Produce json
// @Param register body models.CustomerRegisterRequest true "注册信息"
// @Success 200 {object} models.Response{data=models.CustomerLoginResponse}
// @Router /api/v1/customer/register [post]
func CustomerRegister(c *gin.Context) {
	var req models.CustomerRegisterRequest
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

	// 检查手机号是否已存在（如果提供了手机号）
	if req.PhoneNumber != "" {
		var existingPhone models.Customer
		err = database.DB.Where("phone_number = ?", req.PhoneNumber).First(&existingPhone).Error
		if err == nil {
			utils.Error(c, "手机号已存在")
			return
		} else if err != gorm.ErrRecordNotFound {
			utils.Error(c, "查询手机号失败")
			return
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(c, "密码加密失败")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 创建客户记录
	customer := models.Customer{
		Account:         req.Account,
		PasswordHash:    string(hashedPassword),
		CustomerName:    req.CustomerName,
		ContactMethod:   req.ContactMethod,
		PhoneNumber:     req.PhoneNumber,
		MemberBirthday:  req.MemberBirthday,
		AdditionalInfo1: req.AdditionalInfo1,
		AdditionalInfo2: req.AdditionalInfo2,
		AdditionalInfo3: req.AdditionalInfo3,
		Notes:           req.Notes,
		Status:          "正常",
		AccountType:     "self_registered",
		PasswordStatus:  "set",
	}

	if err := tx.Create(&customer).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建客户失败")
		return
	}

	// 创建财务信息（初始余额为0）
	financialInfo := models.CustomerFinancialInfo{
		CustomerID:        customer.CustomerID,
		InitialRealCharge: 0.00,
		TotalConsumption:  0.00,
		TotalRealCharge:   0.00,
		CurrentBalance:    0.00,
	}
	if err := tx.Create(&financialInfo).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建财务信息失败")
		return
	}

	// 创建偏好设置（默认设置）
	preferences := models.CustomerPreferences{
		CustomerID:             customer.CustomerID,
		ExclusiveDiscountType:  "无折扣",
		ExclusiveDiscountRatio: 0,
		PlatformBoss:           "",
		ExclusiveCS:            "",
	}
	if err := tx.Create(&preferences).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建偏好设置失败")
		return
	}

	// 提交事务
	tx.Commit()

	// 生成JWT令牌
	token, err := utils.GenerateToken(customer.CustomerID, customer.Account, "customer")
	if err != nil {
		log.Printf("[CustomerRegister] 生成访问令牌失败: %v", err)
		utils.Error(c, "生成令牌失败")
		return
	}

	// 重新查询完整信息
	database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, customer.CustomerID)

	response := models.CustomerLoginResponse{
		Token:     token,
		ExpiresIn: 24 * 60 * 60, // 24小时
		User:      &customer,
	}

	utils.SuccessWithMessage(c, "注册成功", response)
}

// CustomerLogin 客户登录
// @Summary 客户登录
// @Description 客户登录接口
// @Tags 客户认证
// @Accept json
// @Produce json
// @Param login body models.CustomerLoginRequest true "登录信息"
// @Success 200 {object} models.Response{data=models.CustomerLoginResponse}
// @Router /api/v1/customer/login [post]
func CustomerLogin(c *gin.Context) {
	var req models.CustomerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.Where("account = ?", req.Account).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, "账户名或密码错误")
		} else {
			utils.Error(c, "登录失败")
		}
		return
	}

	// 检查密码状态
	if customer.PasswordStatus == "unset" {
		utils.Error(c, "账户密码未设置，请联系管理员或使用设置密码功能")
		return
	}

	if customer.PasswordStatus == "need_reset" {
		utils.Error(c, "账户密码需要重置，请使用设置密码功能")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.Password)); err != nil {
		utils.Error(c, "账户名或密码错误")
		return
	}

	// 检查账户状态
	if customer.Status != "正常" {
		utils.Error(c, "账户已被禁用")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(customer.CustomerID, customer.Account, "customer")
	if err != nil {
		log.Printf("[CustomerLogin] 生成访问令牌失败: %v", err)
		utils.Error(c, "生成令牌失败")
		return
	}

	// 更新最后登录信息
	now := time.Now()
	customer.LastLoginAt = &now
	customer.LastLoginIP = c.ClientIP()
	database.DB.Save(&customer)

	// 重新查询完整信息
	database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, customer.CustomerID)

	response := models.CustomerLoginResponse{
		Token:     token,
		ExpiresIn: 24 * 60 * 60, // 24小时
		User:      &customer,
	}

	utils.SuccessWithMessage(c, "登录成功", response)

	log.Printf("[CustomerLogin] 客户登录成功: %s (ID: %d)", customer.Account, customer.CustomerID)
}

// CustomerRefreshToken 刷新客户访问令牌
// @Summary 刷新客户访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 客户认证
// @Accept json
// @Produce json
// @Param refresh body models.CustomerRefreshTokenRequest true "刷新令牌信息"
// @Success 200 {object} models.Response{data=models.CustomerRefreshTokenResponse}
// @Router /api/v1/customer/refresh [post]
func CustomerRefreshToken(c *gin.Context) {
	var req models.CustomerRefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 解析刷新令牌
	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("[CustomerRefreshToken] 解析刷新令牌失败: %v", err)
		utils.ErrorWithCode(c, 401, "刷新令牌无效或已过期")
		return
	}

	// 验证用户角色
	if claims.UserRole != "customer" {
		utils.ErrorWithCode(c, 403, "无效的用户类型")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, claims.MemberID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorWithCode(c, 404, "客户不存在")
		} else {
			utils.ErrorWithCode(c, 500, "查询客户失败")
		}
		return
	}

	// 检查账户状态
	if customer.Status != "正常" {
		utils.ErrorWithCode(c, 403, "账户已被禁用")
		return
	}

	// 生成新的访问令牌
	accessToken, err := utils.GenerateToken(customer.CustomerID, customer.Account, "customer")
	if err != nil {
		log.Printf("[CustomerRefreshToken] 生成访问令牌失败: %v", err)
		utils.ErrorWithCode(c, 500, "生成访问令牌失败")
		return
	}

	// 生成新的刷新令牌
	refreshToken, err := utils.GenerateRefreshToken(customer.CustomerID, customer.Account, "customer")
	if err != nil {
		log.Printf("[CustomerRefreshToken] 生成刷新令牌失败: %v", err)
		utils.ErrorWithCode(c, 500, "生成刷新令牌失败")
		return
	}

	response := models.CustomerRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60, // 24小时
		User:         &customer,
	}

	utils.SuccessWithMessage(c, "令牌刷新成功", response)

	log.Printf("[CustomerRefreshToken] 客户令牌刷新成功: %s (ID: %d)", customer.Account, customer.CustomerID)
}

// GetCustomerProfile 获取客户个人信息
// @Summary 获取客户个人信息
// @Description 获取当前登录客户的个人信息
// @Tags 客户认证
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.Customer}
// @Router /api/v1/customer/profile [get]
func GetCustomerProfile(c *gin.Context) {
	// 从中间件获取客户ID
	customerID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到客户信息")
		return
	}

	var customer models.Customer
	if err := database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, customerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	utils.Success(c, customer)
}

// UpdateCustomerProfile 更新客户个人信息
// @Summary 更新客户个人信息
// @Description 更新当前登录客户的个人信息
// @Tags 客户认证
// @Accept json
// @Produce json
// @Param profile body models.CustomerRegisterRequest true "个人信息"
// @Success 200 {object} models.Response{data=models.Customer}
// @Router /api/v1/customer/profile [put]
func UpdateCustomerProfile(c *gin.Context) {
	// 从中间件获取客户ID
	customerID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到客户信息")
		return
	}

	var req models.CustomerRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateCustomerProfile] 请求参数错误: %v", err)
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.First(&customer, customerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 如果更改了手机号，检查是否重复
	if req.PhoneNumber != "" && req.PhoneNumber != customer.PhoneNumber {
		var existingPhone models.Customer
		err := database.DB.Where("phone_number = ? AND customer_id != ?", req.PhoneNumber, customerID).First(&existingPhone).Error
		if err == nil {
			utils.Error(c, "手机号已存在")
			return
		} else if err != gorm.ErrRecordNotFound {
			utils.Error(c, "查询手机号失败")
			return
		}
	}

	// 更新基本信息（不包括账号和密码）
	customer.CustomerName = req.CustomerName
	customer.ContactMethod = req.ContactMethod
	customer.PhoneNumber = req.PhoneNumber
	customer.MemberBirthday = req.MemberBirthday
	customer.AdditionalInfo1 = req.AdditionalInfo1
	customer.AdditionalInfo2 = req.AdditionalInfo2
	customer.AdditionalInfo3 = req.AdditionalInfo3
	customer.Notes = req.Notes

	if err := database.DB.Save(&customer).Error; err != nil {
		utils.Error(c, "更新客户信息失败")
		return
	}

	// 重新查询完整信息
	database.DB.Preload("FinancialInfo").Preload("Preferences").
		First(&customer, customer.CustomerID)

	utils.SuccessWithMessage(c, "更新成功", customer)
}

// CustomerSetPassword 客户设置密码（首次设置或管理员重置后设置）
// @Summary 客户设置密码
// @Description 为管理员创建的客户账户设置密码，或重置密码后设置新密码
// @Tags 客户认证
// @Accept json
// @Produce json
// @Param password body models.CustomerSetPasswordRequest true "密码信息"
// @Success 200 {object} models.StandardResponse
// @Router /api/v1/customer/set-password [post]
func CustomerSetPassword(c *gin.Context) {
	var req models.CustomerSetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.Where("account = ?", req.Account).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, "账户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 检查密码状态，只有未设置或需要重置的账户才能设置密码
	if customer.PasswordStatus == "set" {
		utils.Error(c, "该账户已设置密码，请使用重置密码功能")
		return
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, "密码加密失败")
		return
	}

	// 更新密码和状态
	updates := map[string]interface{}{
		"password_hash":   string(hashedPassword),
		"password_status": "set",
		"updated_at":      time.Now(),
	}

	if err := database.DB.Model(&customer).Updates(updates).Error; err != nil {
		utils.Error(c, "设置密码失败")
		return
	}

	utils.Success(c, gin.H{
		"success": true,
		"message": "密码设置成功",
	})
}

// CustomerResetPassword 客户重置密码
// @Summary 客户重置密码
// @Description 客户重置自己的登录密码
// @Tags 客户认证
// @Accept json
// @Produce json
// @Param password body models.CustomerResetPasswordRequest true "密码信息"
// @Success 200 {object} models.StandardResponse
// @Router /api/v1/customer/reset-password [post]
func CustomerResetPassword(c *gin.Context) {
	// 从中间件获取客户ID
	customerID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到客户信息")
		return
	}

	var req models.CustomerResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找客户
	var customer models.Customer
	if err := database.DB.First(&customer, customerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "客户不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.OldPassword)); err != nil {
		utils.Error(c, "原密码错误")
		return
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, "密码加密失败")
		return
	}

	// 更新密码
	updates := map[string]interface{}{
		"password_hash": string(hashedPassword),
		"updated_at":    time.Now(),
	}

	if err := database.DB.Model(&customer).Updates(updates).Error; err != nil {
		utils.Error(c, "重置密码失败")
		return
	}

	utils.Success(c, gin.H{
		"success": true,
		"message": "密码重置成功",
	})
}

// GetCustomerBalance 获取客户余额信息
// @Summary 获取客户余额信息
// @Description 获取当前登录客户的余额和财务信息
// @Tags 客户认证
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.CustomerFinancialInfo}
// @Router /api/v1/customer/balance [get]
func GetCustomerBalance(c *gin.Context) {
	// 从中间件获取客户ID
	customerID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到客户信息")
		return
	}

	// 查询财务信息
	var financialInfo models.CustomerFinancialInfo
	if err := database.DB.Where("customer_id = ?", customerID).First(&financialInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "财务信息不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	utils.Success(c, financialInfo)
}
