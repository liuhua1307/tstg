package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderCategory 订单类别表
type OrderCategory struct {
	CategoryID      uint           `json:"category_id" gorm:"primaryKey;column:category_id"`
	CategoryName    string         `json:"category_name" gorm:"uniqueIndex;size:100;not null;comment:订单类别名称"`
	SortOrder       int            `json:"sort_order" gorm:"default:0;comment:排序序号"`
	IsActive        bool           `json:"is_active" gorm:"default:true;comment:是否开启"`
	UsageScenario   string         `json:"usage_scenario" gorm:"size:255;comment:使用场景"`
	CommissionRate  float64        `json:"commission_rate" gorm:"type:decimal(5,2);default:0.00;comment:抽成比例（百分比）"`
	IsParticipating bool           `json:"is_participating" gorm:"default:true;comment:是否参与"`
	IsRequired      bool           `json:"is_required" gorm:"default:false;comment:是否必填"`
	IsAccelerated   bool           `json:"is_accelerated" gorm:"default:false;comment:是否加速"`
	AdditionalInfo  string         `json:"additional_info" gorm:"type:text;comment:附加信息"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (OrderCategory) TableName() string {
	return "order_categories"
}

// PlaymateOrder 陪玩报单基本信息表
type PlaymateOrder struct {
	OrderID               uint      `json:"order_id" gorm:"primaryKey;column:order_id"`
	ReporterID            uint      `json:"reporter_id" gorm:"not null;comment:报单人ID"`
	CustomerID            uint      `json:"customer_id" gorm:"not null;comment:客户ID"`
	OrderCategoryID       uint      `json:"order_category_id" gorm:"not null;comment:订单类别ID"`
	ProjectCategory       string    `json:"project_category" gorm:"size:100;not null;comment:项目分类"`
	StartTime             time.Time `json:"start_time" gorm:"not null;comment:开始时间"`
	EndTime               time.Time `json:"end_time" gorm:"not null;comment:结束时间"`
	DurationHours         float64   `json:"duration_hours" gorm:"type:decimal(5,2);not null;comment:陪玩时长（小时）"`
	ServiceAdditionalInfo string    `json:"service_additional_info" gorm:"type:text;comment:服务附加说明"`
	InternalNotes         string    `json:"internal_notes" gorm:"type:text;comment:内部备注"`
	OrderNotes            string    `json:"order_notes" gorm:"type:text;comment:订单备注"`
	ReportTime            time.Time `json:"report_time" gorm:"comment:报单时间"`
	CustomerName          string    `json:"customer_name" gorm:"comment:客户名称"`
	//新增字段
	UseBalancePayment bool           `json:"use_balance_payment" gorm:"comment:是否用余额结算"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Reporter    *InternalMember   `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	Customer    *Customer         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Category    *OrderCategory    `json:"category,omitempty" gorm:"foreignKey:OrderCategoryID"`
	Pricing     *OrderPricing     `json:"pricing,omitempty" gorm:"foreignKey:OrderID"`
	Workflow    *OrderWorkflow    `json:"workflow,omitempty" gorm:"foreignKey:OrderID"`
	PaymentInfo *OrderPaymentInfo `json:"payment_info,omitempty" gorm:"foreignKey:OrderID"`
	Images      []OrderImages     `json:"images,omitempty" gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (PlaymateOrder) TableName() string {
	return "playmate_orders"
}

// OrderPricing 订单价格信息表
type OrderPricing struct {
	PricingID         uint      `json:"pricing_id" gorm:"primaryKey;column:pricing_id"`
	OrderID           uint      `json:"order_id" gorm:"uniqueIndex;not null;comment:订单ID"`
	UnitPrice         float64   `json:"unit_price" gorm:"type:decimal(10,2);not null;comment:单价（元/小时）"`
	TotalPrice        float64   `json:"total_price" gorm:"type:decimal(10,2);not null;comment:订单总价"`
	DiscountAmount    float64   `json:"discount_amount" gorm:"type:decimal(10,2);default:0.00;comment:折扣总价"`
	FinalPrice        float64   `json:"final_price" gorm:"type:decimal(10,2);not null;comment:最终结算价格"`
	ExclusiveDiscount bool      `json:"exclusive_discount" gorm:"default:false;comment:是否专属折扣"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// 关联关系
	Order *PlaymateOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (OrderPricing) TableName() string {
	return "order_pricing"
}

// OrderWorkflow 订单工作流状态表
type OrderWorkflow struct {
	WorkflowID      uint       `json:"workflow_id" gorm:"primaryKey;column:workflow_id"`
	OrderID         uint       `json:"order_id" gorm:"uniqueIndex;not null;comment:订单ID"`
	OrderStatus     string     `json:"order_status" gorm:"type:enum('待处理','驳回','已确认','已退回');default:'待处理';comment:订单状态"`
	ApproverID      *uint      `json:"approver_id" gorm:"comment:审批人ID"`
	ApprovalTime    *time.Time `json:"approval_time" gorm:"comment:审批时间"`
	RejectionReason string     `json:"rejection_reason" gorm:"type:text;comment:驳回原因"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// 关联关系
	Order    *PlaymateOrder  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Approver *InternalMember `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
}

// TableName 指定表名
func (OrderWorkflow) TableName() string {
	return "order_workflow"
}

// OrderPaymentInfo 订单支付信息表
type OrderPaymentInfo struct {
	PaymentID     uint       `json:"payment_id" gorm:"autoIncrement:false;column:payment_id"` // 移除主键，禁用自增
	OrderID       uint       `json:"order_id" gorm:"primaryKey;not null;comment:订单ID"`      // 改为主键
	TransactionID string     `json:"transaction_id" gorm:"type:text;comment:付款流水号"`      // 修复：从 payment_transaction_id 改为 transaction_id
	PaymentMethod string     `json:"payment_method" gorm:"size:50;comment:付款方式"`
	PaymentTime   *time.Time `json:"payment_time" gorm:"comment:付款时间"`
	PaymentAmount float64    `json:"payment_amount" gorm:"type:decimal(10,2);comment:付款金额"`
	PaymentStatus string     `json:"payment_status" gorm:"type:enum('待付款','已付款','付款失败','已退款');default:'待付款';comment:付款状态"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// 关联关系
	Order *PlaymateOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (OrderPaymentInfo) TableName() string {
	return "order_payment_info"
}

// OrderImages 订单图片表
type OrderImages struct {
	ImageID   uint      `json:"image_id" gorm:"primaryKey;column:image_id"`
	OrderID   uint      `json:"order_id" gorm:"not null;comment:关联的订单ID"`
	ImageURL  string    `json:"image_url" gorm:"size:255;not null;comment:图片存储URL"`
	ImageType string    `json:"image_type" gorm:"type:enum('结算截图','聊天记录','游戏截图','其他');default:'其他';comment:图片类型"`
	UploadAt  time.Time `json:"upload_at" gorm:"comment:上传时间"`
	Notes     string    `json:"notes" gorm:"type:text;comment:图片备注"`

	// 关联关系
	Order *PlaymateOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (OrderImages) TableName() string {
	return "order_images"
}

// 请求结构
type OrderCreateRequest struct {
	CustomerID            uint      `json:"customer_id" binding:"required"`
	OrderCategoryID       uint      `json:"order_category_id" binding:"required"`
	ProjectCategory       string    `json:"project_category"`
	StartTime             time.Time `json:"start_time" binding:"required"` // 修复：从 string 改为 time.Time
	EndTime               time.Time `json:"end_time" binding:"required"`   // 修复：从 string 改为 time.Time
	DurationHours         float64   `json:"duration_hours" binding:"required"`
	UnitPrice             float64   `json:"unit_price" binding:"required"`
	ServiceAdditionalInfo string    `json:"service_additional_info"`
	InternalNotes         string    `json:"internal_notes"`
	OrderNotes            string    `json:"order_notes"`
	UseBalancePayment     bool      `json:"use_balance_payment"` // 是否使用余额支付
	ExclusiveDiscount     bool      `json:"exclusive_discount"`
}

type OrderUpdateStatusRequest struct {
	OrderStatus     string `json:"order_status" binding:"required"`
	RejectionReason string `json:"rejection_reason"`
}

type OrderCategoryCreateRequest struct {
	CategoryName    string  `json:"category_name" binding:"required"`
	SortOrder       int     `json:"sort_order"`
	UsageScenario   string  `json:"usage_scenario"`
	CommissionRate  float64 `json:"commission_rate"`
	IsParticipating bool    `json:"is_participating"`
	IsRequired      bool    `json:"is_required"`
	IsAccelerated   bool    `json:"is_accelerated"`
	AdditionalInfo  string  `json:"additional_info"`
}

// OrderApprovalHistory 订单审批操作历史表
type OrderApprovalHistory struct {
	ActionID     uint           `json:"action_id" gorm:"primaryKey;column:action_id"`
	OrderID      uint           `json:"order_id" gorm:"not null;comment:订单ID"`
	OperatorID   uint           `json:"operator_id" gorm:"not null;comment:操作人ID"`
	OperatorName string         `json:"operator_name" gorm:"size:100;comment:操作人姓名"`
	Action       string         `json:"action" gorm:"type:enum('approve','reject','status_change');comment:操作类型"`
	FromStatus   string         `json:"from_status" gorm:"size:50;comment:原状态"`
	ToStatus     string         `json:"to_status" gorm:"size:50;comment:新状态"`
	Reason       string         `json:"reason" gorm:"type:text;comment:操作原因"`
	Notes        string         `json:"notes" gorm:"type:text;comment:备注"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Order    *PlaymateOrder  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Operator *InternalMember `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

// TableName 指定表名
func (OrderApprovalHistory) TableName() string {
	return "order_approval_history"
}

// 订单审批相关请求结构
type OrderApprovalFilterRequest struct {
	PageRequest
	OrderID    string  `json:"order_id" form:"order_id"`
	Status     string  `json:"status" form:"status"`
	ReporterID uint    `json:"reporter_id" form:"reporter_id"`
	CustomerID uint    `json:"customer_id" form:"customer_id"`
	Category   string  `json:"category" form:"category"`
	StartDate  string  `json:"start_date" form:"start_date"`
	EndDate    string  `json:"end_date" form:"end_date"`
	DateType   string  `json:"date_type" form:"date_type"` // submit, approve, settle
	MinAmount  float64 `json:"min_amount" form:"min_amount"`
	MaxAmount  float64 `json:"max_amount" form:"max_amount"`
	SortBy     string  `json:"sort_by" form:"sort_by"`
	SortOrder  string  `json:"sort_order" form:"sort_order"`
}

type ApproveOrderRequest struct {
	Notes string `json:"notes"`
}

type RejectOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
	Notes  string `json:"notes"`
}

type BatchApprovalRequest struct {
	OrderIDs []string `json:"order_ids" binding:"required"`
	Action   string   `json:"action" binding:"required"`
	Reason   string   `json:"reason"`
	Notes    string   `json:"notes"`
}

type UpdateOrderStatusRequestV2 struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
	Notes  string `json:"notes"`
}

type GetStatisticsRequest struct {
	OrderID    string  `json:"order_id" form:"order_id"`
	Status     string  `json:"status" form:"status"`
	ReporterID uint    `json:"reporter_id" form:"reporter_id"`
	CustomerID uint    `json:"customer_id" form:"customer_id"`
	Category   string  `json:"category" form:"category"`
	StartDate  string  `json:"start_date" form:"start_date"`
	EndDate    string  `json:"end_date" form:"end_date"`
	DateType   string  `json:"date_type" form:"date_type"`
	MinAmount  float64 `json:"min_amount" form:"min_amount"`
	MaxAmount  float64 `json:"max_amount" form:"max_amount"`
}

// 订单统计相关结构体
type GetOrderStatsReq struct {
	ReporterID      *uint  `json:"reporter_id" form:"reporter_id"`
	CustomerID      *uint  `json:"customer_id" form:"customer_id"`
	OrderCategoryID *uint  `json:"order_category_id" form:"order_category_id"`
	OrderStatus     string `json:"order_status" form:"order_status"`
	StartDate       string `json:"start_date" form:"start_date"`
	EndDate         string `json:"end_date" form:"end_date"`
}

type GetOrderStatsResData struct {
	TotalCount         int     `json:"totalCount"`
	TotalDurationHours float64 `json:"totalDurationHours"`
	TotalAmount        float64 `json:"totalAmount"`
	TotalCommission    float64 `json:"totalCommission"`
}

type GetOperationHistoryRequest struct {
	PageRequest
	OrderID    string `json:"order_id" form:"order_id"`
	OperatorID uint   `json:"operator_id" form:"operator_id"`
	Action     string `json:"action" form:"action"`
	StartDate  string `json:"start_date" form:"start_date"`
	EndDate    string `json:"end_date" form:"end_date"`
}

type BatchExportRequest struct {
	OrderIDs []string `json:"order_ids" binding:"required"`
	Format   string   `json:"format" binding:"required"`
	Fields   []string `json:"fields"`
}

// 响应结构
type OrderStatistics struct {
	TotalCount         int64            `json:"total_count"`
	TotalHours         float64          `json:"total_hours"`
	TotalAmount        float64          `json:"total_amount"`
	TotalCommission    float64          `json:"total_commission"`
	AveragePrice       float64          `json:"average_price"`
	StatusDistribution map[string]int64 `json:"status_distribution"`
	DailyTrend         []DailyTrendItem `json:"daily_trend"`
}

type DailyTrendItem struct {
	Date   string  `json:"date"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
}

type BatchApprovalResponse struct {
	Success      bool                   `json:"success"`
	SuccessCount int                    `json:"success_count"`
	FailureCount int                    `json:"failure_count"`
	Failures     []BatchApprovalFailure `json:"failures"`
}

type BatchApprovalFailure struct {
	OrderID string `json:"order_id"`
	Error   string `json:"error"`
}

type BatchExportResponse struct {
	DownloadURL string `json:"download_url"`
	Filename    string `json:"filename"`
}

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
