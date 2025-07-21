package controllers

import (
	"strconv"
	"tangsong-esports/database"
	"tangsong-esports/models"
	"tangsong-esports/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetOrderCategories 获取订单类别列表
// @Summary 获取订单类别列表
// @Description 获取所有订单类别，支持筛选
// @Tags 订单类别管理
// @Accept json
// @Produce json
// @Param is_active query bool false "是否启用"
// @Success 200 {object} models.Response{data=[]models.OrderCategory}
// @Router /api/v1/order-categories [get]
func GetOrderCategories(c *gin.Context) {
	query := database.DB.Model(&models.OrderCategory{})

	// 筛选条件
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			query = query.Where("is_active = ?", isActive)
		}
	}

	var categories []models.OrderCategory
	if err := query.Order("sort_order ASC, category_id ASC").Find(&categories).Error; err != nil {
		utils.Error(c, "查询失败")
		return
	}

	utils.Success(c, categories)
}

// CreateOrderCategory 创建订单类别
// @Summary 创建订单类别
// @Description 创建新的订单类别
// @Tags 订单类别管理
// @Accept json
// @Produce json
// @Param category body models.OrderCategoryCreateRequest true "订单类别信息"
// @Success 200 {object} models.Response{data=models.OrderCategory}
// @Router /api/v1/order-categories [post]
func CreateOrderCategory(c *gin.Context) {
	var req models.OrderCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 检查类别名是否已存在
	var existingCategory models.OrderCategory
	err := database.DB.Where("category_name = ?", req.CategoryName).First(&existingCategory).Error
	if err == nil {
		utils.Error(c, "类别名称已存在")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.Error(c, "查询类别失败")
		return
	}

	// 创建订单类别
	category := models.OrderCategory{
		CategoryName:    req.CategoryName,
		SortOrder:       req.SortOrder,
		IsActive:        true,
		UsageScenario:   req.UsageScenario,
		CommissionRate:  req.CommissionRate,
		IsParticipating: req.IsParticipating,
		IsRequired:      req.IsRequired,
		IsAccelerated:   req.IsAccelerated,
		AdditionalInfo:  req.AdditionalInfo,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		utils.Error(c, "创建订单类别失败")
		return
	}

	utils.SuccessWithMessage(c, "创建订单类别成功", category)
}

// UpdateOrderCategory 更新订单类别
// @Summary 更新订单类别
// @Description 更新订单类别信息
// @Tags 订单类别管理
// @Accept json
// @Produce json
// @Param id path int true "类别ID"
// @Param category body models.OrderCategoryCreateRequest true "订单类别信息"
// @Success 200 {object} models.Response{data=models.OrderCategory}
// @Router /api/v1/order-categories/{id} [put]
func UpdateOrderCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的类别ID")
		return
	}

	var req models.OrderCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找订单类别
	var category models.OrderCategory
	if err := database.DB.First(&category, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "订单类别不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 更新信息
	category.CategoryName = req.CategoryName
	category.SortOrder = req.SortOrder
	category.UsageScenario = req.UsageScenario
	category.CommissionRate = req.CommissionRate
	category.IsParticipating = req.IsParticipating
	category.IsRequired = req.IsRequired
	category.IsAccelerated = req.IsAccelerated
	category.AdditionalInfo = req.AdditionalInfo

	if err := database.DB.Save(&category).Error; err != nil {
		utils.Error(c, "更新订单类别失败")
		return
	}

	utils.SuccessWithMessage(c, "更新订单类别成功", category)
}

// DeleteOrderCategory 删除订单类别
// @Summary 删除订单类别
// @Description 软删除订单类别
// @Tags 订单类别管理
// @Accept json
// @Produce json
// @Param id path int true "类别ID"
// @Success 200 {object} models.Response
// @Router /api/v1/order-categories/{id} [delete]
func DeleteOrderCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的类别ID")
		return
	}

	// 查找订单类别
	var category models.OrderCategory
	if err := database.DB.First(&category, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "订单类别不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 软删除
	if err := database.DB.Delete(&category).Error; err != nil {
		utils.Error(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除订单类别成功", nil)
}

// GetOrders 获取订单列表
// @Summary 获取陪玩订单列表
// @Description 分页获取陪玩订单列表，支持搜索和筛选
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param reporter_id query int false "报单人ID"
// @Param customer_id query int false "客户ID"
// @Param game query string false "游戏名称"
// @Param order_status query string false "订单状态"
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/orders [get]
func GetOrders(c *gin.Context) {
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
	query := database.DB.Model(&models.PlaymateOrder{})

	// 搜索条件
	if reporterIDStr := c.Query("reporter_id"); reporterIDStr != "" {
		if reporterID, err := strconv.ParseUint(reporterIDStr, 10, 32); err == nil {
			query = query.Where("reporter_id = ?", uint(reporterID))
		}
	}
	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		if customerID, err := strconv.ParseUint(customerIDStr, 10, 32); err == nil {
			query = query.Where("customer_id = ?", uint(customerID))
		}
	}
	if game := c.Query("game"); game != "" {
		query = query.Where("game LIKE ?", "%"+game+"%")
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("DATE(start_time) >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("DATE(end_time) <= ?", endDate)
	}

	// 关联查询订单状态
	if orderStatus := c.Query("order_status"); orderStatus != "" {
		query = query.Joins("JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id").
			Where("order_workflow.order_status = ?", orderStatus)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var orders []models.PlaymateOrder
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Reporter").Preload("Customer").Preload("Category").
		Preload("Pricing").Preload("Workflow").Preload("PaymentInfo").
		Order("report_time DESC").Offset(offset).Limit(req.PageSize).Find(&orders).Error

	if err != nil {
		utils.Error(c, "查询失败")
		return
	}

	response := models.PageResponse{
		List:     orders,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.Success(c, response)
}

// CreateOrder 创建订单
// @Summary 创建陪玩订单
// @Description 创建新的陪玩订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param order body models.OrderCreateRequest true "订单信息"
// @Success 200 {object} models.Response{data=models.PlaymateOrder}
// @Router /api/v1/orders [post]
func CreateOrder(c *gin.Context) {
	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 获取报单人ID
	reporterID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到报单人信息")
		return
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		utils.Error(c, "开始时间格式错误")
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		utils.Error(c, "结束时间格式错误")
		return
	}

	// 验证客户和订单类别存在
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		utils.Error(c, "客户不存在")
		return
	}

	var category models.OrderCategory
	if err := database.DB.First(&category, req.OrderCategoryID).Error; err != nil {
		utils.Error(c, "订单类别不存在")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 创建订单基本信息
	order := models.PlaymateOrder{
		ReporterID:            reporterID.(uint),
		CustomerID:            req.CustomerID,
		OrderCategoryID:       req.OrderCategoryID,
		Game:                  req.Game,
		ProjectCategory:       req.ProjectCategory,
		PlaymateLevel:         req.PlaymateLevel,
		StartTime:             startTime,
		EndTime:               endTime,
		DurationHours:         req.DurationHours,
		IsTeammate:            req.IsTeammate,
		Mode:                  req.Mode,
		ServiceAdditionalInfo: req.ServiceAdditionalInfo,
		InternalNotes:         req.InternalNotes,
		OrderNotes:            req.OrderNotes,
		PlatformOwner:         req.PlatformOwner,
		ReportTime:            time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建订单失败")
		return
	}

	// 计算价格
	totalPrice := req.UnitPrice * req.DurationHours
	finalPrice := totalPrice
	discountAmount := 0.0
	if req.ExclusiveDiscount {
		// 这里可以根据客户的专属折扣计算
		discountAmount = totalPrice * 0.1 // 示例：10%折扣
		finalPrice = totalPrice - discountAmount
	}

	// 创建价格信息
	pricing := models.OrderPricing{
		OrderID:           order.OrderID,
		UnitPrice:         req.UnitPrice,
		TotalPrice:        totalPrice,
		DiscountAmount:    discountAmount,
		FinalPrice:        finalPrice,
		ExclusiveDiscount: req.ExclusiveDiscount,
	}
	if err := tx.Create(&pricing).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建价格信息失败")
		return
	}

	// 创建工作流状态
	workflow := models.OrderWorkflow{
		OrderID:          order.OrderID,
		OrderStatus:      "待处理",
		SettlementStatus: "未结算",
	}
	if err := tx.Create(&workflow).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建工作流失败")
		return
	}

	// 创建支付信息
	paymentInfo := models.OrderPaymentInfo{
		OrderID:       order.OrderID,
		PaymentAmount: finalPrice,
		PaymentStatus: "待付款",
	}
	if err := tx.Create(&paymentInfo).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "创建支付信息失败")
		return
	}

	// 提交事务
	tx.Commit()

	// 重新查询完整信息
	database.DB.Preload("Reporter").Preload("Customer").Preload("Category").
		Preload("Pricing").Preload("Workflow").Preload("PaymentInfo").
		First(&order, order.OrderID)

	utils.SuccessWithMessage(c, "创建订单成功", order)
}

// UpdateOrder 更新订单
// @Summary 更新陪玩订单
// @Description 更新陪玩订单信息
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param order body models.OrderCreateRequest true "订单信息"
// @Success 200 {object} models.Response{data=models.PlaymateOrder}
// @Router /api/v1/orders/{id} [put]
func UpdateOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的订单ID")
		return
	}

	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找订单
	var order models.PlaymateOrder
	if err := database.DB.First(&order, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "订单不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 检查订单状态是否允许修改
	var workflow models.OrderWorkflow
	if err := database.DB.Where("order_id = ?", order.OrderID).First(&workflow).Error; err == nil {
		if workflow.OrderStatus != "待处理" {
			utils.Error(c, "只有待处理状态的订单才能修改")
			return
		}
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		utils.Error(c, "开始时间格式错误")
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		utils.Error(c, "结束时间格式错误")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 更新订单基本信息
	order.Game = req.Game
	order.ProjectCategory = req.ProjectCategory
	order.PlaymateLevel = req.PlaymateLevel
	order.StartTime = startTime
	order.EndTime = endTime
	order.DurationHours = req.DurationHours
	order.IsTeammate = req.IsTeammate
	order.Mode = req.Mode
	order.ServiceAdditionalInfo = req.ServiceAdditionalInfo
	order.InternalNotes = req.InternalNotes
	order.OrderNotes = req.OrderNotes
	order.PlatformOwner = req.PlatformOwner

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "更新订单失败")
		return
	}

	// 更新价格信息
	var pricing models.OrderPricing
	if err := tx.Where("order_id = ?", order.OrderID).First(&pricing).Error; err == nil {
		totalPrice := req.UnitPrice * req.DurationHours
		finalPrice := totalPrice
		discountAmount := 0.0
		if req.ExclusiveDiscount {
			discountAmount = totalPrice * 0.1
			finalPrice = totalPrice - discountAmount
		}

		pricing.UnitPrice = req.UnitPrice
		pricing.TotalPrice = totalPrice
		pricing.DiscountAmount = discountAmount
		pricing.FinalPrice = finalPrice
		pricing.ExclusiveDiscount = req.ExclusiveDiscount
		tx.Save(&pricing)
	}

	// 提交事务
	tx.Commit()

	// 重新查询完整信息
	database.DB.Preload("Reporter").Preload("Customer").Preload("Category").
		Preload("Pricing").Preload("Workflow").Preload("PaymentInfo").
		First(&order, order.OrderID)

	utils.SuccessWithMessage(c, "更新订单成功", order)
}

// GetOrderByID 根据ID获取订单
// @Summary 获取订单详情
// @Description 根据ID获取陪玩订单详细信息
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} models.Response{data=models.PlaymateOrder}
// @Router /api/v1/orders/{id} [get]
func GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的订单ID")
		return
	}

	var order models.PlaymateOrder
	if err := database.DB.Preload("Reporter").Preload("Customer").Preload("Category").
		Preload("Pricing").Preload("Workflow").Preload("PaymentInfo").Preload("Images").
		First(&order, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "订单不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	utils.Success(c, order)
}

// UpdateOrderStatus 更新订单状态
// @Summary 更新订单状态
// @Description 更新陪玩订单状态
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param status body models.OrderUpdateStatusRequest true "状态信息"
// @Success 200 {object} models.Response{data=models.OrderWorkflow}
// @Router /api/v1/orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的订单ID")
		return
	}

	var req models.OrderUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 获取审批人ID
	approverID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到审批人信息")
		return
	}

	// 查找订单工作流
	var workflow models.OrderWorkflow
	if err := database.DB.Where("order_id = ?", uint(id)).First(&workflow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "订单不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 更新状态
	workflow.OrderStatus = req.OrderStatus
	workflow.ApproverID = &[]uint{approverID.(uint)}[0]
	now := time.Now()
	workflow.ApprovalTime = &now

	if req.OrderStatus == "驳回" {
		workflow.RejectionReason = req.RejectionReason
	}

	if req.OrderStatus == "已结算" {
		workflow.SettlementStatus = "已结算"
		workflow.SettlementTime = &now
	}

	if err := database.DB.Save(&workflow).Error; err != nil {
		utils.Error(c, "更新状态失败")
		return
	}

	// 重新查询完整信息
	database.DB.Preload("Order").Preload("Approver").First(&workflow, workflow.WorkflowID)

	utils.SuccessWithMessage(c, "更新状态成功", workflow)
}

// UploadOrderImages 上传订单图片
// @Summary 上传订单图片
// @Description 为订单上传图片
// @Tags 订单管理
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "订单ID"
// @Param image_type formData string false "图片类型"
// @Param notes formData string false "图片备注"
// @Success 200 {object} models.Response{data=models.OrderImages}
// @Router /api/v1/orders/{id}/images [post]
func UploadOrderImages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, "无效的订单ID")
		return
	}

	// 验证订单是否存在
	var order models.PlaymateOrder
	if err := database.DB.First(&order, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "订单不存在")
		} else {
			utils.Error(c, "查询订单失败")
		}
		return
	}

	// 获取表单数据
	imageType := c.PostForm("image_type")
	notes := c.PostForm("notes")

	// 这里应该实现文件上传逻辑
	// 暂时返回模拟数据
	orderImage := models.OrderImages{
		OrderID:   uint(id),
		ImageURL:  "/uploads/order_images/example.jpg", // 实际应该是上传后的文件路径
		ImageType: imageType,
		Notes:     notes,
		UploadAt:  time.Now(),
	}

	if err := database.DB.Create(&orderImage).Error; err != nil {
		utils.Error(c, "保存图片记录失败")
		return
	}

	utils.SuccessWithMessage(c, "上传图片成功", orderImage)
}
