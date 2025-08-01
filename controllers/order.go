package controllers

import (
	"fmt"
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
	//判断当前用户身份
	memberID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到报单人信息")
		return
	}
	var role string
	isAdmin := true
	query := database.DB.Model(&models.PlaymateOrder{}).Where("playmate_orders.reporter_id = ?", memberID)
	// 检查权限是不是超级管理员、管理员
	database.DB.Model(&models.InternalMember{}).Select("user_role").Where("member_id = ?", memberID).First(&role)
	if role != "超级管理员" && role != "管理员" {
		isAdmin = false
	}
	// 设置默认值

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询，如果是管理员，则查询所有订单，否则查询当前用户报单的订单
	if isAdmin {
		query = database.DB.Model(&models.PlaymateOrder{})
	} else {
		query = database.DB.Model(&models.PlaymateOrder{}).Where("playmate_orders.reporter_id = ?", memberID)
	}

	// 搜索条件
	if reporterIDStr := c.Query("reporter_id"); reporterIDStr != "" {
		if reporterID, err := strconv.ParseUint(reporterIDStr, 10, 32); err == nil {
			query = query.Where("playmate_orders.reporter_id = ?", uint(reporterID))
		}
	}
	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		if customerID, err := strconv.ParseUint(customerIDStr, 10, 32); err == nil {
			query = query.Where("playmate_orders.customer_id = ?", uint(customerID))
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("DATE(playmate_orders.start_time) >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("DATE(playmate_orders.end_time) <= ?", endDate)
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
	err := query.
		Select("playmate_orders.*").
		Joins("LEFT JOIN customers ON playmate_orders.customer_id = customers.customer_id").
		Preload("Reporter").Preload("Customer").Preload("Category").
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
		fmt.Println("请求参数错误", err)
		utils.Error(c, "请求参数错误")
		return
	}

	// 获取报单人ID
	reporterID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到报单人信息")
		return
	}

	// 验证客户和订单类别存在
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		utils.Error(c, "客户不存在")
		return
	}
	fmt.Println(customer)
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
		ProjectCategory:       req.ProjectCategory,
		StartTime:             req.StartTime,
		EndTime:               req.EndTime,
		DurationHours:         req.DurationHours,
		ServiceAdditionalInfo: req.ServiceAdditionalInfo,
		InternalNotes:         req.InternalNotes,
		OrderNotes:            req.OrderNotes,
		CustomerName:          customer.CustomerName,
		UseBalancePayment:     req.UseBalancePayment,
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
		OrderID:     order.OrderID,
		OrderStatus: "待处理",
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

	// 开启事务
	tx := database.DB.Begin()

	// 更新订单基本信息
	order.ProjectCategory = req.ProjectCategory
	order.StartTime = req.StartTime
	order.EndTime = req.EndTime
	order.DurationHours = req.DurationHours
	order.ServiceAdditionalInfo = req.ServiceAdditionalInfo
	order.InternalNotes = req.InternalNotes
	order.OrderNotes = req.OrderNotes

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
	if err := database.DB.
		Preload("Reporter").Preload("Customer").Preload("Category").
		Preload("Pricing").Preload("Workflow").Preload("PaymentInfo").Preload("Images").
		First(&order, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "订单不存在")
		} else {
			utils.Error(c, "查询失败")
		}
		return
	}

	// 手动查询客户信息并设置customerName
	var customer models.Customer
	if err := database.DB.Where("customer_id = ? AND deleted_at IS NULL", order.CustomerID).First(&customer).Error; err == nil {
		order.CustomerName = customer.CustomerName
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

// GetOrderStats 获取订单统计数据
// @Summary 获取订单统计数据
// @Description 根据条件统计订单数量、时长、金额和佣金
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param reporter_id query int false "报单人ID"
// @Param customer_id query int false "客户ID"
// @Param order_category_id query int false "订单类别ID"
// @Param game query string false "游戏名称"
// @Param order_status query string false "订单状态"
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Success 200 {object} models.Response{data=models.GetOrderStatsResData}
// @Router /api/v1/orders/stats [get]
func GetOrderStats(c *gin.Context) {
	var req models.GetOrderStatsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 构建基础查询
	query := database.DB.Model(&models.PlaymateOrder{}).
		Joins("LEFT JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id").
		Joins("LEFT JOIN order_pricing ON playmate_orders.order_id = order_pricing.order_id").
		Joins("LEFT JOIN order_categories ON playmate_orders.order_category_id = order_categories.category_id")

	// 应用筛选条件
	if req.ReporterID != nil {
		query = query.Where("playmate_orders.reporter_id = ?", *req.ReporterID)
	}
	if req.CustomerID != nil {
		query = query.Where("playmate_orders.customer_id = ?", *req.CustomerID)
	}
	if req.OrderCategoryID != nil {
		query = query.Where("playmate_orders.order_category_id = ?", *req.OrderCategoryID)
	}
	if req.OrderStatus != "" {
		query = query.Where("order_workflow.order_status = ?", req.OrderStatus)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(playmate_orders.start_time) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(playmate_orders.end_time) <= ?", req.EndDate)
	}

	// 执行统计查询
	var result struct {
		TotalCount         int64   `json:"total_count"`
		TotalDurationHours float64 `json:"total_duration_hours"`
		TotalAmount        float64 `json:"total_amount"`
		TotalCommission    float64 `json:"total_commission"`
	}

	err := query.Select(`
		COUNT(DISTINCT playmate_orders.order_id) as total_count,
		COALESCE(SUM(playmate_orders.duration_hours), 0) as total_duration_hours,
		COALESCE(SUM(order_pricing.final_price), 0) as total_amount,
		COALESCE(SUM(order_pricing.final_price * order_categories.commission_rate), 0) as total_commission
	`).Scan(&result).Error

	if err != nil {
		utils.Error(c, "统计查询失败")
		return
	}

	// 构建响应数据
	statsData := models.GetOrderStatsResData{
		TotalCount:         int(result.TotalCount),
		TotalDurationHours: result.TotalDurationHours,
		TotalAmount:        result.TotalAmount,
		TotalCommission:    result.TotalCommission,
	}

	utils.Success(c, statsData)
}

// ===================订单审批相关接口===================

// GetPendingOrders 获取待审批订单列表
// @Summary 获取待审批订单列表
// @Description 获取所有待处理状态的订单列表
// @Tags 订单审批
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/order-approval/pending [get]
func GetPendingOrders(c *gin.Context) {
	var req models.OrderApprovalFilterRequest
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

	// 构建查询，只查询待处理状态的订单
	query := database.DB.Model(&models.PlaymateOrder{}).
		Joins("JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id").
		Where("order_workflow.order_status = ?", "待处理")

	// 应用筛选条件
	query = applyOrderFilters(query, req)

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var orders []models.PlaymateOrder
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Select("playmate_orders.* ").
		Joins("LEFT JOIN customers ON playmate_orders.customer_id = customers.customer_id").
		Preload("Reporter").Preload("Customer").Preload("Category").
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

// GetApprovalOrders 获取审批订单列表（包含所有状态）
// @Summary 获取审批订单列表
// @Description 获取所有状态的订单列表，支持筛选
// @Tags 订单审批
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/order-approval [get]
func GetApprovalOrders(c *gin.Context) {
	var req models.OrderApprovalFilterRequest
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
	query := database.DB.Model(&models.PlaymateOrder{}).
		Joins("JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id")

	// 应用筛选条件
	query = applyOrderFilters(query, req)

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var orders []models.PlaymateOrder
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Select("playmate_orders.*").
		Joins("LEFT JOIN customers ON playmate_orders.customer_id = customers.customer_id").
		Preload("Reporter").Preload("Customer").Preload("Category").
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

// ApproveOrder 审批通过订单
// @Summary 审批通过订单
// @Description 将订单状态更新为已确认，同时从客户余额中扣除订单金额并更新支付状态为已付款
// @Tags 订单审批
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param data body models.ApproveOrderRequest true "审批信息"
// @Success 200 {object} models.Response{data=models.StandardResponse}
// @Router /api/v1/order-approval/{id}/approve [post]
func ApproveOrder(c *gin.Context) {
	orderID := c.Param("id")
	id, err := strconv.ParseUint(orderID, 10, 32)
	if err != nil {
		utils.Error(c, "无效的订单ID")
		return
	}

	var req models.ApproveOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 获取操作人ID
	operatorID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到操作人信息")
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

	// 检查当前状态
	if workflow.OrderStatus != "待处理" {
		utils.Error(c, "只有待处理状态的订单才能审批")
		return
	}

	// 获取操作人信息
	var operator models.InternalMember
	if err := database.DB.First(&operator, operatorID).Error; err != nil {
		utils.Error(c, "获取操作人信息失败")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 记录旧状态
	oldStatus := workflow.OrderStatus

	// 更新订单状态
	workflow.OrderStatus = "已确认"
	workflow.ApproverID = &[]uint{operatorID.(uint)}[0]
	now := time.Now()
	workflow.ApprovalTime = &now

	if err := tx.Save(&workflow).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "更新订单状态失败")
		return
	}

	// 记录操作历史
	history := models.OrderApprovalHistory{
		OrderID:      uint(id),
		OperatorID:   operatorID.(uint),
		OperatorName: operator.Name,
		Action:       "approve",
		FromStatus:   oldStatus,
		ToStatus:     "已确认",
		Notes:        req.Notes,
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "记录操作历史失败")
		return
	}

	// === 新增扣款逻辑 ===
	// 获取订单信息（包括客户ID和价格信息）
	var order models.PlaymateOrder
	if err := tx.Preload("Pricing").First(&order, uint(id)).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "获取订单信息失败")
		return
	}
	// 检查订单是否允许使用余额支付
	if order.UseBalancePayment == false {
		// 更新订单支付状态
		var paymentInfo models.OrderPaymentInfo
		if err := tx.Where("order_id = ?", uint(id)).First(&paymentInfo).Error; err != nil {
			tx.Rollback()
			utils.Error(c, "获取支付信息失败")
			return
		}
		paymentInfo.PaymentStatus = "已付款"
		paymentInfo.PaymentMethod = "直接支付"
		paymentTime := time.Now()
		paymentInfo.PaymentTime = &paymentTime
		if err := tx.Save(&paymentInfo).Error; err != nil {
			tx.Rollback()
			utils.Error(c, "更新支付状态失败")
			return
		}
		// 提交事务
		tx.Commit()
		response := models.StandardResponse{
			Success: true,
			Message: fmt.Sprintf("订单审批成功,未走余额支付流程，订单ID: %d", id),
		}

		utils.SuccessWithMessage(c, "订单审批成功", response)
		return
	}
	// 获取客户财务信息
	var customerFinancial models.CustomerFinancialInfo
	if err := tx.Where("customer_id = ?", order.CustomerID).First(&customerFinancial).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			utils.Error(c, "客户财务信息不存在，请先为客户充值")
			return
		} else {
			tx.Rollback()
			utils.Error(c, "获取客户财务信息失败")
			return
		}
	}

	// 获取客户偏好设置（包括折扣信息）
	var customerPreferences models.CustomerPreferences
	if err := tx.Where("customer_id = ?", order.CustomerID).First(&customerPreferences).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "获取客户偏好设置失败")
		return
	}

	// 计算订单金额和折扣
	originalAmount := order.Pricing.FinalPrice
	discountRate := 0.0
	actualAmount := originalAmount

	// 直接使用客户设置的折扣比例（0-100）计算实际扣款金额
	if customerPreferences.ExclusiveDiscountRatio >= 0 && customerPreferences.ExclusiveDiscountRatio <= 100 {
		discountRate = float64(customerPreferences.ExclusiveDiscountRatio) / 100.0
		actualAmount = originalAmount * (1 - discountRate)
	} else {
		// 如果折扣比例无效，按原价处理
		discountRate = 0.0
		actualAmount = originalAmount
	}

	// 检查余额是否足够（使用折扣后的金额）
	if customerFinancial.CurrentBalance < actualAmount {
		tx.Rollback()
		utils.Error(c, fmt.Sprintf("客户余额不足，当前余额：%.2f，订单原价：%.2f，折扣后金额：%.2f（%.0f%%折扣）",
			customerFinancial.CurrentBalance, originalAmount, actualAmount, discountRate*100))
		return
	}

	// 扣除客户余额（使用折扣后的金额）
	customerFinancial.CurrentBalance -= actualAmount
	customerFinancial.TotalConsumption += actualAmount
	if err := tx.Save(&customerFinancial).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "扣款失败")
		return
	}

	// 更新订单支付状态
	var paymentInfo models.OrderPaymentInfo
	if err := tx.Where("order_id = ?", uint(id)).First(&paymentInfo).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "获取支付信息失败")
		return
	}

	paymentInfo.PaymentStatus = "已付款"
	paymentInfo.PaymentMethod = "余额支付"
	paymentTime := time.Now()
	paymentInfo.PaymentTime = &paymentTime
	if err := tx.Save(&paymentInfo).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "更新支付状态失败")
		return
	}

	// 提交事务
	tx.Commit()

	response := models.StandardResponse{
		Success: true,
		Message: fmt.Sprintf("订单审批成功！原价：%.2f元，折扣：%.0f%%，实际扣款：%.2f元",
			originalAmount, discountRate*100, actualAmount),
	}

	utils.SuccessWithMessage(c, "订单审批成功", response)
}

// RejectOrder 驳回订单
// @Summary 驳回订单
// @Description 将订单状态更新为驳回
// @Tags 订单审批
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param data body models.RejectOrderRequest true "驳回信息"
// @Success 200 {object} models.Response{data=models.StandardResponse}
// @Router /api/v1/order-approval/{id}/reject [post]
func RejectOrder(c *gin.Context) {
	orderID := c.Param("id")
	id, err := strconv.ParseUint(orderID, 10, 32)
	if err != nil {
		utils.Error(c, "无效的订单ID")
		return
	}

	var req models.RejectOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 获取操作人ID
	operatorID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到操作人信息")
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

	// 检查当前状态
	if workflow.OrderStatus != "待处理" {
		utils.Error(c, "只有待处理状态的订单才能驳回")
		return
	}

	// 获取操作人信息
	var operator models.InternalMember
	if err := database.DB.First(&operator, operatorID).Error; err != nil {
		utils.Error(c, "获取操作人信息失败")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 记录旧状态
	oldStatus := workflow.OrderStatus

	// 更新订单状态
	workflow.OrderStatus = "驳回"
	workflow.RejectionReason = req.Reason
	workflow.ApproverID = &[]uint{operatorID.(uint)}[0]
	now := time.Now()
	workflow.ApprovalTime = &now

	if err := tx.Save(&workflow).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "更新订单状态失败")
		return
	}

	// 记录操作历史
	history := models.OrderApprovalHistory{
		OrderID:      uint(id),
		OperatorID:   operatorID.(uint),
		OperatorName: operator.Name,
		Action:       "reject",
		FromStatus:   oldStatus,
		ToStatus:     "驳回",
		Reason:       req.Reason,
		Notes:        req.Notes,
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "记录操作历史失败")
		return
	}

	tx.Commit()

	response := models.StandardResponse{
		Success: true,
		Message: "订单驳回成功",
	}

	utils.SuccessWithMessage(c, "订单驳回成功", response)
}

// BatchApproval 批量审批
// @Summary 批量审批订单
// @Description 批量审批或驳回订单
// @Tags 订单审批
// @Accept json
// @Produce json
// @Param data body models.BatchApprovalRequest true "批量审批信息"
// @Success 200 {object} models.Response{data=models.BatchApprovalResponse}
// @Router /api/v1/order-approval/batch [post]
func BatchApproval(c *gin.Context) {
	var req models.BatchApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 获取操作人ID
	operatorID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到操作人信息")
		return
	}

	// 获取操作人信息
	var operator models.InternalMember
	if err := database.DB.First(&operator, operatorID).Error; err != nil {
		utils.Error(c, "获取操作人信息失败")
		return
	}

	var failures []models.BatchApprovalFailure
	successCount := 0

	// 逐个处理订单
	for _, orderIDStr := range req.OrderIDs {
		orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
		if err != nil {
			failures = append(failures, models.BatchApprovalFailure{
				OrderID: orderIDStr,
				Error:   "无效的订单ID",
			})
			continue
		}

		// 查找订单工作流
		var workflow models.OrderWorkflow
		if err := database.DB.Where("order_id = ?", uint(orderID)).First(&workflow).Error; err != nil {
			failures = append(failures, models.BatchApprovalFailure{
				OrderID: orderIDStr,
				Error:   "订单不存在",
			})
			continue
		}

		// 检查当前状态
		if workflow.OrderStatus != "待处理" {
			failures = append(failures, models.BatchApprovalFailure{
				OrderID: orderIDStr,
				Error:   "只有待处理状态的订单才能审批",
			})
			continue
		}

		// 开启事务
		tx := database.DB.Begin()

		// 记录旧状态
		oldStatus := workflow.OrderStatus
		var newStatus string
		var action string

		switch req.Action {
		case "approve":
			newStatus = "已确认"
			action = "approve"
		case "reject":
			newStatus = "驳回"
			action = "reject"
			workflow.RejectionReason = req.Reason
		default:
			failures = append(failures, models.BatchApprovalFailure{
				OrderID: orderIDStr,
				Error:   "无效的操作类型",
			})
			continue
		}

		// 更新订单状态
		workflow.OrderStatus = newStatus
		workflow.ApproverID = &[]uint{operatorID.(uint)}[0]
		now := time.Now()
		workflow.ApprovalTime = &now

		if err := tx.Save(&workflow).Error; err != nil {
			tx.Rollback()
			failures = append(failures, models.BatchApprovalFailure{
				OrderID: orderIDStr,
				Error:   "更新订单状态失败",
			})
			continue
		}

		// 记录操作历史
		history := models.OrderApprovalHistory{
			OrderID:      uint(orderID),
			OperatorID:   operatorID.(uint),
			OperatorName: operator.Name,
			Action:       action,
			FromStatus:   oldStatus,
			ToStatus:     newStatus,
			Reason:       req.Reason,
			Notes:        req.Notes,
		}

		if err := tx.Create(&history).Error; err != nil {
			tx.Rollback()
			failures = append(failures, models.BatchApprovalFailure{
				OrderID: orderIDStr,
				Error:   "记录操作历史失败",
			})
			continue
		}

		tx.Commit()
		successCount++
	}

	response := models.BatchApprovalResponse{
		Success:      len(failures) == 0,
		SuccessCount: successCount,
		FailureCount: len(failures),
		Failures:     failures,
	}

	utils.Success(c, response)
}

// UpdateOrderStatusV2 更新订单状态（审批模块使用）
// @Summary 更新订单状态
// @Description 更新订单状态
// @Tags 订单审批
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param data body models.UpdateOrderStatusRequestV2 true "状态更新信息"
// @Success 200 {object} models.Response{data=models.StandardResponse}
// @Router /api/v1/order-approval/{id}/status [patch]
func UpdateOrderStatusV2(c *gin.Context) {
	orderID := c.Param("id")
	id, err := strconv.ParseUint(orderID, 10, 32)
	if err != nil {
		utils.Error(c, "无效的订单ID")
		return
	}

	var req models.UpdateOrderStatusRequestV2
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 获取操作人ID
	operatorID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到操作人信息")
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

	// 获取操作人信息
	var operator models.InternalMember
	if err := database.DB.First(&operator, operatorID).Error; err != nil {
		utils.Error(c, "获取操作人信息失败")
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 记录旧状态
	oldStatus := workflow.OrderStatus

	// 更新订单状态
	workflow.OrderStatus = req.Status
	workflow.ApproverID = &[]uint{operatorID.(uint)}[0]
	now := time.Now()
	workflow.ApprovalTime = &now

	if req.Status == "驳回" && req.Reason != "" {
		workflow.RejectionReason = req.Reason
	}

	if err := tx.Save(&workflow).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "更新订单状态失败")
		return
	}

	// 记录操作历史
	history := models.OrderApprovalHistory{
		OrderID:      uint(id),
		OperatorID:   operatorID.(uint),
		OperatorName: operator.Name,
		Action:       "status_change",
		FromStatus:   oldStatus,
		ToStatus:     req.Status,
		Reason:       req.Reason,
		Notes:        req.Notes,
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		utils.Error(c, "记录操作历史失败")
		return
	}

	tx.Commit()

	response := models.StandardResponse{
		Success: true,
		Message: "状态更新成功",
	}

	utils.SuccessWithMessage(c, "状态更新成功", response)
}

// GetStatistics 获取统计数据
// @Summary 获取订单统计数据
// @Description 获取订单的统计信息
// @Tags 订单审批
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.OrderStatistics}
// @Router /api/v1/order-approval/statistics [get]
func GetStatistics(c *gin.Context) {
	var req models.GetStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 构建基础查询
	query := database.DB.Model(&models.PlaymateOrder{}).
		Joins("JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id").
		Joins("JOIN order_pricing ON playmate_orders.order_id = order_pricing.order_id")

	// 应用筛选条件
	filterReq := models.OrderApprovalFilterRequest{
		OrderID:    req.OrderID,
		Status:     req.Status,
		ReporterID: req.ReporterID,
		CustomerID: req.CustomerID,
		Category:   req.Category,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		DateType:   req.DateType,
		MinAmount:  req.MinAmount,
		MaxAmount:  req.MaxAmount,
	}
	query = applyOrderFilters(query, filterReq)

	// 计算基础统计
	var stats struct {
		TotalCount  int64   `json:"total_count"`
		TotalHours  float64 `json:"total_hours"`
		TotalAmount float64 `json:"total_amount"`
	}

	err := query.Select("COUNT(*) as total_count, SUM(duration_hours) as total_hours, SUM(order_pricing.final_price) as total_amount").
		Scan(&stats).Error
	if err != nil {
		utils.Error(c, "查询统计数据失败")
		return
	}

	// 计算平均价格和佣金
	averagePrice := float64(0)
	if stats.TotalCount > 0 {
		averagePrice = stats.TotalAmount / float64(stats.TotalCount)
	}
	totalCommission := stats.TotalAmount * 0.1 // 假设10%的佣金率

	// 获取状态分布
	var statusDistribution []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	statusQuery := database.DB.Model(&models.PlaymateOrder{}).
		Joins("JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id").
		Select("order_workflow.order_status as status, COUNT(*) as count").
		Group("order_workflow.order_status")

	// 应用相同的筛选条件（除了状态筛选）
	filterReqForStatus := filterReq
	filterReqForStatus.Status = ""
	statusQuery = applyOrderFilters(statusQuery, filterReqForStatus)

	if err := statusQuery.Scan(&statusDistribution).Error; err != nil {
		utils.Error(c, "查询状态分布失败")
		return
	}

	// 转换状态分布格式
	statusMap := make(map[string]int64)
	for _, item := range statusDistribution {
		statusMap[item.Status] = item.Count
	}

	// 获取每日趋势（最近30天）
	var dailyTrend []models.DailyTrendItem
	dailyQuery := database.DB.Model(&models.PlaymateOrder{}).
		Joins("JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id").
		Joins("JOIN order_pricing ON playmate_orders.order_id = order_pricing.order_id").
		Select("DATE(report_time) as date, COUNT(*) as count, SUM(order_pricing.final_price) as amount").
		Where("report_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)").
		Group("DATE(report_time)").
		Order("date DESC")

	// 应用筛选条件
	dailyQuery = applyOrderFilters(dailyQuery, filterReq)

	if err := dailyQuery.Scan(&dailyTrend).Error; err != nil {
		utils.Error(c, "查询每日趋势失败")
		return
	}

	response := models.OrderStatistics{
		TotalCount:         stats.TotalCount,
		TotalHours:         stats.TotalHours,
		TotalAmount:        stats.TotalAmount,
		TotalCommission:    totalCommission,
		AveragePrice:       averagePrice,
		StatusDistribution: statusMap,
		DailyTrend:         dailyTrend,
	}

	utils.Success(c, response)
}

// GetOperationHistory 获取操作历史
// @Summary 获取操作历史
// @Description 获取订单的操作历史记录
// @Tags 订单审批
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/order-approval/history [get]
func GetOperationHistory(c *gin.Context) {
	var req models.GetOperationHistoryRequest
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
	query := database.DB.Model(&models.OrderApprovalHistory{})

	// 应用筛选条件
	if req.OrderID != "" {
		orderID, err := strconv.ParseUint(req.OrderID, 10, 32)
		if err == nil {
			query = query.Where("order_id = ?", uint(orderID))
		}
	}
	if req.OperatorID > 0 {
		query = query.Where("operator_id = ?", req.OperatorID)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(created_at) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(created_at) <= ?", req.EndDate)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var history []models.OrderApprovalHistory
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Order").Preload("Operator").
		Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&history).Error

	if err != nil {
		utils.Error(c, "查询失败")
		return
	}

	response := models.PageResponse{
		List:     history,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.Success(c, response)
}

// BatchExport 批量导出
// @Summary 批量导出订单
// @Description 批量导出订单数据
// @Tags 订单审批
// @Accept json
// @Produce json
// @Param data body models.BatchExportRequest true "导出请求"
// @Success 200 {object} models.Response{data=models.BatchExportResponse}
// @Router /api/v1/order-approval/export [post]
func BatchExport(c *gin.Context) {
	var req models.BatchExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 验证格式
	if req.Format != "excel" && req.Format != "pdf" {
		utils.Error(c, "不支持的导出格式")
		return
	}

	// 查询要导出的订单
	var orderIDs []uint
	for _, idStr := range req.OrderIDs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			orderIDs = append(orderIDs, uint(id))
		}
	}

	var orders []models.PlaymateOrder
	err := database.DB.Where("order_id IN ?", orderIDs).
		Preload("Reporter").Preload("Customer").Preload("Category").
		Preload("Pricing").Preload("Workflow").Preload("PaymentInfo").
		Find(&orders).Error

	if err != nil {
		utils.Error(c, "查询订单失败")
		return
	}

	// 这里应该实现实际的导出逻辑
	// 暂时返回模拟数据
	filename := fmt.Sprintf("orders_export_%d.%s", time.Now().Unix(), req.Format)
	downloadURL := fmt.Sprintf("/downloads/exports/%s", filename)

	response := models.BatchExportResponse{
		DownloadURL: downloadURL,
		Filename:    filename,
	}

	utils.SuccessWithMessage(c, "导出任务已创建", response)
}

// applyOrderFilters 应用订单筛选条件的辅助函数
func applyOrderFilters(query *gorm.DB, req models.OrderApprovalFilterRequest) *gorm.DB {
	if req.OrderID != "" {
		if orderID, err := strconv.ParseUint(req.OrderID, 10, 32); err == nil {
			query = query.Where("playmate_orders.order_id = ?", uint(orderID))
		}
	}
	if req.Status != "" && req.Status != "all" {
		query = query.Where("order_workflow.order_status = ?", req.Status)
	}

	if req.ReporterID > 0 {
		query = query.Where("playmate_orders.reporter_id = ?", req.ReporterID)
	}
	if req.CustomerID > 0 {
		query = query.Where("playmate_orders.customer_id = ?", req.CustomerID)
	}

	if req.Category != "" {
		query = query.Joins("JOIN order_categories ON playmate_orders.order_category_id = order_categories.category_id").
			Where("order_categories.category_name LIKE ?", "%"+req.Category+"%")
	}

	// 日期筛选
	if req.StartDate != "" {
		switch req.DateType {
		case "approve":
			query = query.Where("DATE(order_workflow.approval_time) >= ?", req.StartDate)
		default: // submit
			query = query.Where("DATE(playmate_orders.report_time) >= ?", req.StartDate)
		}
	}
	if req.EndDate != "" {
		switch req.DateType {
		case "approve":
			query = query.Where("DATE(order_workflow.approval_time) <= ?", req.EndDate)
		default: // submit
			query = query.Where("DATE(playmate_orders.report_time) <= ?", req.EndDate)
		}
	}

	// 金额筛选
	if req.MinAmount > 0 {
		query = query.Joins("JOIN order_pricing ON playmate_orders.order_id = order_pricing.order_id").
			Where("order_pricing.final_price >= ?", req.MinAmount)
	}
	if req.MaxAmount > 0 {
		query = query.Joins("JOIN order_pricing ON playmate_orders.order_id = order_pricing.order_id").
			Where("order_pricing.final_price <= ?", req.MaxAmount)
	}

	// 排序
	if req.SortBy != "" {
		sortField := "playmate_orders.report_time" // 默认排序字段
		switch req.SortBy {
		case "order_id":
			sortField = "playmate_orders.order_id"
		case "start_time":
			sortField = "playmate_orders.start_time"
		case "end_time":
			sortField = "playmate_orders.end_time"
		case "approval_time":
			sortField = "order_workflow.approval_time"
		}

		sortOrder := "DESC"
		if req.SortOrder == "asc" {
			sortOrder = "ASC"
		}

		query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
	}

	return query
}

// GetCustomerOrders 客户查看自己的订单
// @Summary 客户查看自己的订单列表
// @Description 客户查看自己下过的所有订单，支持分页和状态筛选
// @Tags 客户订单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param order_status query string false "订单状态筛选" Enums(待处理,已确认,驳回,已退回)
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/customer/orders [get]
func GetCustomerOrders(c *gin.Context) {
	// 获取当前客户ID
	customerID, exists := c.Get("member_id")
	if !exists {
		utils.Error(c, "未找到客户信息")
		return
	}

	// 检查用户角色是否为客户
	userRole, roleExists := c.Get("user_role")
	if !roleExists || userRole != "customer" {
		utils.ErrorWithCode(c, 403, "无权限访问")
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

	// 构建查询，只查询当前客户的订单
	query := database.DB.Model(&models.PlaymateOrder{}).Where("playmate_orders.customer_id = ?", customerID)

	// 搜索条件
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("DATE(playmate_orders.start_time) >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("DATE(playmate_orders.end_time) <= ?", endDate)
	}

	// 关联查询订单状态
	if orderStatus := c.Query("order_status"); orderStatus != "" {
		query = query.Joins("JOIN order_workflow ON playmate_orders.order_id = order_workflow.order_id").
			Where("order_workflow.order_status = ?", orderStatus)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询，加载相关信息
	var orders []models.PlaymateOrder
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Select("playmate_orders.*").
		Joins("LEFT JOIN customers ON playmate_orders.customer_id = customers.customer_id").
		Preload("Reporter").Preload("Customer").Preload("Category").
		Preload("Pricing").Preload("Workflow").Preload("PaymentInfo").
		Order("report_time DESC").Offset(offset).Limit(req.PageSize).Find(&orders).Error

	if err != nil {
		utils.Error(c, "查询失败")
		return
	}

	// 构建返回数据
	response := models.PageResponse{
		List:     orders,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.Success(c, response)
}
