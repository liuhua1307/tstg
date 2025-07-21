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
	OrderID               uint           `json:"order_id" gorm:"primaryKey;column:order_id"`
	ReporterID            uint           `json:"reporter_id" gorm:"not null;comment:报单人ID"`
	CustomerID            uint           `json:"customer_id" gorm:"not null;comment:客户ID"`
	OrderCategoryID       uint           `json:"order_category_id" gorm:"not null;comment:订单类别ID"`
	Game                  string         `json:"game" gorm:"size:100;not null;comment:游戏名称"`
	ProjectCategory       string         `json:"project_category" gorm:"size:100;not null;comment:项目分类"`
	PlaymateLevel         string         `json:"playmate_level" gorm:"size:50;comment:陪玩等级"`
	StartTime             time.Time      `json:"start_time" gorm:"not null;comment:开始时间"`
	EndTime               time.Time      `json:"end_time" gorm:"not null;comment:结束时间"`
	DurationHours         float64        `json:"duration_hours" gorm:"type:decimal(5,2);not null;comment:陪玩时长（小时）"`
	IsTeammate            bool           `json:"is_teammate" gorm:"default:false;comment:是否队友"`
	Mode                  string         `json:"mode" gorm:"size:50;comment:报单模式"`
	ServiceAdditionalInfo string         `json:"service_additional_info" gorm:"type:text;comment:服务附加说明"`
	InternalNotes         string         `json:"internal_notes" gorm:"type:text;comment:内部备注"`
	OrderNotes            string         `json:"order_notes" gorm:"type:text;comment:订单备注"`
	PlatformOwner         string         `json:"platform_owner" gorm:"size:100;comment:所属平台老板"`
	ReportTime            time.Time      `json:"report_time" gorm:"comment:报单时间"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`

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
	WorkflowID       uint       `json:"workflow_id" gorm:"primaryKey;column:workflow_id"`
	OrderID          uint       `json:"order_id" gorm:"uniqueIndex;not null;comment:订单ID"`
	OrderStatus      string     `json:"order_status" gorm:"type:enum('待处理','驳回','已确认','已结算','已完成','已退回');default:'待处理';comment:订单状态"`
	SettlementStatus string     `json:"settlement_status" gorm:"type:enum('未结算','已结算');default:'未结算';comment:结算状态"`
	ApproverID       *uint      `json:"approver_id" gorm:"comment:审批人ID"`
	ApprovalTime     *time.Time `json:"approval_time" gorm:"comment:审批时间"`
	RejectionReason  string     `json:"rejection_reason" gorm:"type:text;comment:驳回原因"`
	SettlementTime   *time.Time `json:"settlement_time" gorm:"comment:结算时间"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

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
	PaymentID     uint       `json:"payment_id" gorm:"primaryKey;column:payment_id"`                    // 修复：从 payment_info_id 改为 payment_id
	OrderID       uint       `json:"order_id" gorm:"uniqueIndex;not null;comment:订单ID"`
	TransactionID string     `json:"transaction_id" gorm:"type:text;comment:付款流水号"`                    // 修复：从 payment_transaction_id 改为 transaction_id
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
	Game                  string    `json:"game" binding:"required"`
	ProjectCategory       string    `json:"project_category" binding:"required"`
	PlaymateLevel         string    `json:"playmate_level"`
	StartTime             time.Time `json:"start_time" binding:"required"`          // 修复：从 string 改为 time.Time
	EndTime               time.Time `json:"end_time" binding:"required"`            // 修复：从 string 改为 time.Time
	DurationHours         float64   `json:"duration_hours" binding:"required"`
	UnitPrice             float64   `json:"unit_price" binding:"required"`
	IsTeammate            bool      `json:"is_teammate"`
	Mode                  string    `json:"mode"`
	ServiceAdditionalInfo string    `json:"service_additional_info"`
	InternalNotes         string    `json:"internal_notes"`
	OrderNotes            string    `json:"order_notes"`
	PlatformOwner         string    `json:"platform_owner"`
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