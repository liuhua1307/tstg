package controllers

import (
	"tangsong-esports/database"
	"tangsong-esports/models"
	"tangsong-esports/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetConfigs 获取系统配置
// @Summary 获取系统配置
// @Description 获取所有系统配置
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param is_active query bool false "是否启用"
// @Success 200 {object} models.Response{data=[]models.SystemConfig}
// @Router /api/v1/configs [get]
func GetConfigs(c *gin.Context) {
	query := database.DB.Model(&models.SystemConfig{})

	// 筛选条件
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActiveStr == "true" {
			query = query.Where("is_active = ?", true)
		} else if isActiveStr == "false" {
			query = query.Where("is_active = ?", false)
		}
	}

	var configs []models.SystemConfig
	if err := query.Order("config_key ASC").Find(&configs).Error; err != nil {
		utils.Error(c, "查询配置失败")
		return
	}

	utils.Success(c, configs)
}

// UpdateConfig 更新系统配置
// @Summary 更新系统配置
// @Description 更新指定配置项的值
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param key path string true "配置键"
// @Param config body map[string]interface{} true "配置信息"
// @Success 200 {object} models.Response{data=models.SystemConfig}
// @Router /api/v1/configs/{key} [put]
func UpdateConfig(c *gin.Context) {
	configKey := c.Param("key")
	if configKey == "" {
		utils.Error(c, "配置键不能为空")
		return
	}

	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找配置
	var config models.SystemConfig
	if err := database.DB.Where("config_key = ?", configKey).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "配置不存在")
		} else {
			utils.Error(c, "查询配置失败")
		}
		return
	}

	// 更新配置值
	if value, exists := requestData["config_value"]; exists {
		if valueStr, ok := value.(string); ok {
			config.ConfigValue = valueStr
		}
	}

	if description, exists := requestData["config_description"]; exists {
		if descStr, ok := description.(string); ok {
			config.ConfigDescription = descStr
		}
	}

	if isActive, exists := requestData["is_active"]; exists {
		if activeVal, ok := isActive.(bool); ok {
			config.IsActive = activeVal
		}
	}

	// 保存配置
	if err := database.DB.Save(&config).Error; err != nil {
		utils.Error(c, "更新配置失败")
		return
	}

	// 记录操作日志
	operatorID, _ := c.Get("member_id")
	logOperation(operatorID.(uint), "更新", "系统配置", "更新系统配置："+configKey, configKey, "系统配置", c.ClientIP(), c.GetHeader("User-Agent"))

	utils.SuccessWithMessage(c, "更新配置成功", config)
}

// GetOperationLogs 获取操作日志
// @Summary 获取操作日志
// @Description 分页获取操作日志，支持筛选
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param operator_id query int false "操作人ID"
// @Param operation_type query string false "操作类型"
// @Param operation_module query string false "操作模块"
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Success 200 {object} models.Response{data=models.PageResponse}
// @Router /api/v1/logs [get]
func GetOperationLogs(c *gin.Context) {
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
	query := database.DB.Model(&models.OperationLog{})

	// 搜索条件
	if operatorIDStr := c.Query("operator_id"); operatorIDStr != "" {
		query = query.Where("operator_id = ?", operatorIDStr)
	}
	if operationType := c.Query("operation_type"); operationType != "" {
		query = query.Where("operation_type LIKE ?", "%"+operationType+"%")
	}
	if operationModule := c.Query("operation_module"); operationModule != "" {
		query = query.Where("operation_module = ?", operationModule)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("DATE(operation_time) >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("DATE(operation_time) <= ?", endDate)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var logs []models.OperationLog
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Operator").Order("operation_time DESC").
		Offset(offset).Limit(req.PageSize).Find(&logs).Error

	if err != nil {
		utils.Error(c, "查询日志失败")
		return
	}

	response := models.PageResponse{
		List:     logs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	utils.Success(c, response)
}

// logOperation 记录操作日志的辅助函数
func logOperation(operatorID uint, operationType, operationModule, description, targetID, targetType, ipAddress, userAgent string) {
	operationLog := models.OperationLog{
		OperatorID:           operatorID,
		OperationType:        operationType,
		OperationModule:      operationModule,
		OperationDescription: description,
		TargetID:             targetID,
		TargetType:           targetType,
		IPAddress:            ipAddress,
		UserAgent:            userAgent,
	}

	// 异步记录日志，不影响主业务流程
	go func() {
		database.DB.Create(&operationLog)
	}()
}
