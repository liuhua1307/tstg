package models

import (
	"time"

	"gorm.io/gorm"
)

// Customer 客户基本信息表
type Customer struct {
	CustomerID      uint           `json:"customer_id" gorm:"primaryKey;column:customer_id"`
	Account         string         `json:"account" gorm:"uniqueIndex;size:50;not null;comment:登录账号"`
	CustomerName    string         `json:"customer_name" gorm:"size:100;not null;comment:客户名称/昵称"`
	ContactMethod   string         `json:"contact_method" gorm:"size:100;comment:联系方式（微信/QQ/手机）"`
	PhoneNumber     string         `json:"phone_number" gorm:"uniqueIndex;size:20;comment:手机号码"`
	MemberBirthday  *time.Time     `json:"member_birthday" gorm:"type:date;comment:会员生日"`
	RoomCode        string         `json:"room_code" gorm:"size:50;comment:房间码"`
	AdditionalInfo1 string         `json:"additional_info1" gorm:"type:text;comment:附加信息1"`
	AdditionalInfo2 string         `json:"additional_info2" gorm:"type:text;comment:附加信息2"`
	AdditionalInfo3 string         `json:"additional_info3" gorm:"type:text;comment:附加信息3"`
	Notes           string         `json:"notes" gorm:"type:text;comment:备注"`
	Status          string         `json:"status" gorm:"type:enum('正常','禁用','过期');default:'正常';comment:客户账户状态"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	FinancialInfo   *CustomerFinancialInfo    `json:"financial_info,omitempty" gorm:"foreignKey:CustomerID"`
	Preferences     *CustomerPreferences      `json:"preferences,omitempty" gorm:"foreignKey:CustomerID"`
	RechargeHistory []CustomerRechargeHistory `json:"recharge_history,omitempty" gorm:"foreignKey:CustomerID"`
}

// TableName 指定表名
func (Customer) TableName() string {
	return "customers"
}

// CustomerFinancialInfo 客户财务信息表
type CustomerFinancialInfo struct {
	FinancialID       uint      `json:"financial_id" gorm:"primaryKey;column:financial_id"`
	CustomerID        uint      `json:"customer_id" gorm:"uniqueIndex;not null;comment:客户ID"`
	InitialRealCharge float64   `json:"initial_real_charge" gorm:"type:decimal(10,2);default:0.00;comment:初始实充金额"`
	TotalConsumption  float64   `json:"total_consumption" gorm:"type:decimal(10,2);default:0.00;comment:历史消费总额"`
	TotalRealCharge   float64   `json:"total_real_charge" gorm:"type:decimal(10,2);default:0.00;comment:历史实际充值总额"`
	CurrentBalance    float64   `json:"current_balance" gorm:"type:decimal(10,2);default:0.00;comment:当前可用余额"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// 关联关系
	Customer *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}

// TableName 指定表名
func (CustomerFinancialInfo) TableName() string {
	return "customer_financial_info"
}

// CustomerPreferences 客户服务偏好表
type CustomerPreferences struct {
	PreferenceID          uint      `json:"preference_id" gorm:"primaryKey;column:preference_id"`
	CustomerID            uint      `json:"customer_id" gorm:"uniqueIndex;not null;comment:客户ID"`
	ExclusiveDiscountType string    `json:"exclusive_discount_type" gorm:"type:enum('无折扣','固定折扣');default:'无折扣';comment:专属折扣类型"`
	PlatformBoss          string    `json:"platform_boss" gorm:"size:100;comment:所属平台老板"`
	ExclusiveCS           string    `json:"exclusive_cs" gorm:"size:100;comment:专属服务客服"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// 关联关系
	Customer *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}

// TableName 指定表名
func (CustomerPreferences) TableName() string {
	return "customer_preferences"
}

// CustomerRechargeHistory 客户充值记录表
type CustomerRechargeHistory struct {
	RechargeID          uint      `json:"recharge_id" gorm:"primaryKey;column:recharge_id"`
	CustomerID          uint      `json:"customer_id" gorm:"not null;comment:客户ID"`
	RealChargeAmount    float64   `json:"real_charge_amount" gorm:"type:decimal(10,2);not null;comment:实充金额"`
	GiftAmount          float64   `json:"gift_amount" gorm:"type:decimal(10,2);default:0.00;comment:赠送金额"`
	TotalRechargeAmount float64   `json:"total_recharge_amount" gorm:"type:decimal(10,2);not null;comment:本次充值总额"`
	PaymentMethod       string    `json:"payment_method" gorm:"type:enum('微信','支付宝','银行转账','平台','内部','其他');not null;comment:付款方式"`
	TransactionID       string    `json:"transaction_id" gorm:"type:text;comment:收款单号/交易ID"`
	RechargeAt          time.Time `json:"recharge_at" gorm:"comment:充值时间"`
	Notes               string    `json:"notes" gorm:"type:text;comment:备注"`
	OperatorID          uint      `json:"operator_id" gorm:"not null;comment:操作员ID"`

	// 关联关系
	Customer *Customer       `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Operator *InternalMember `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

// TableName 指定表名
func (CustomerRechargeHistory) TableName() string {
	return "customer_recharge_history"
}

// 请求结构
type CustomerCreateRequest struct {
	Account               string     `json:"account" binding:"required"`
	CustomerName          string     `json:"customer_name" binding:"required"`
	ContactMethod         string     `json:"contact_method"`
	PhoneNumber           string     `json:"phone_number"`
	MemberBirthday        *time.Time `json:"member_birthday"`        // 修复：改为 *time.Time 类型
	RoomCode              string     `json:"room_code"`
	AdditionalInfo1       string     `json:"additional_info1"`
	AdditionalInfo2       string     `json:"additional_info2"`
	AdditionalInfo3       string     `json:"additional_info3"`
	Notes                 string     `json:"notes"`
	InitialRealCharge     float64    `json:"initial_real_charge"`
	ExclusiveDiscountType string     `json:"exclusive_discount_type"`
	PlatformBoss          string     `json:"platform_boss"`
	ExclusiveCS           string     `json:"exclusive_cs"`
}

type CustomerRechargeRequest struct {
	CustomerID       uint    `json:"customer_id"`
	RealChargeAmount float64 `json:"real_charge_amount" binding:"required"`
	GiftAmount       float64 `json:"gift_amount"`
	PaymentMethod    string  `json:"payment_method" binding:"required"`
	TransactionID    string  `json:"transaction_id"`
	Notes            string  `json:"notes"`
}