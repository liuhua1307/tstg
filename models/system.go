package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	ConfigID          uint           `json:"config_id" gorm:"primaryKey;column:config_id"`
	ConfigKey         string         `json:"config_key" gorm:"uniqueIndex;size:100;not null;comment:配置键"`
	ConfigValue       string         `json:"config_value" gorm:"type:text;comment:配置值"`
	ConfigDescription string         `json:"config_description" gorm:"size:255;comment:配置描述"`
	IsActive          bool           `json:"is_active" gorm:"default:true;comment:是否启用"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// OperationLog 操作日志表
type OperationLog struct {
	LogID                uint      `json:"log_id" gorm:"primaryKey;column:log_id"`
	OperatorID           uint      `json:"operator_id" gorm:"not null;comment:操作人ID"`
	OperationType        string    `json:"operation_type" gorm:"size:50;not null;comment:操作类型"`
	OperationModule      string    `json:"operation_module" gorm:"size:50;comment:操作模块"`
	OperationDescription string    `json:"operation_description" gorm:"type:text;comment:操作描述"`
	TargetID             string    `json:"target_id" gorm:"size:36;comment:操作目标ID"`
	TargetType           string    `json:"target_type" gorm:"size:50;comment:操作目标类型"`
	OldData              string    `json:"old_data" gorm:"type:json;comment:操作前数据"`
	NewData              string    `json:"new_data" gorm:"type:json;comment:操作后数据"`
	IPAddress            string    `json:"ip_address" gorm:"size:45;comment:IP地址"`
	UserAgent            string    `json:"user_agent" gorm:"type:text;comment:用户代理"`
	OperationTime        time.Time `json:"operation_time" gorm:"comment:操作时间"`

	// 关联关系
	Operator *InternalMember `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 分页请求
type PageRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 常量定义
const (
	StatusSuccess   = 200
	StatusError     = 500
	StatusForbidden = 403
	StatusNotFound  = 404
)
