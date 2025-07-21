package models

import (
	"time"

	"gorm.io/gorm"
)

// InternalMember 内部成员基本信息表
type InternalMember struct {
	MemberID     uint           `json:"member_id" gorm:"primaryKey;column:member_id"`
	Account      string         `json:"account" gorm:"uniqueIndex;size:50;not null;comment:登录账号"`
	PasswordHash string         `json:"-" gorm:"size:255;not null;comment:密码哈希值"`
	Name         string         `json:"name" gorm:"size:50;not null;comment:姓名"`
	PhoneNumber  string         `json:"phone_number" gorm:"uniqueIndex;size:20;comment:手机号码"`
	Department   string         `json:"department" gorm:"size:50;comment:所属部门"`
	UserRole     string         `json:"user_role" gorm:"size:50;not null;comment:用户角色"`
	Status       string         `json:"status" gorm:"type:enum('正常','禁用');default:'正常';comment:账户状态"`
	IsEnabled    bool           `json:"is_enabled" gorm:"default:true;comment:是否开启"`
	Notes        string         `json:"notes" gorm:"type:text;comment:备注"`
	LastLoginAt  *time.Time     `json:"last_login_at" gorm:"comment:最近登录时间"`
	LastLoginIP  string         `json:"last_login_ip" gorm:"size:45;comment:最近登录IP"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Permissions       *MemberPermissions       `json:"permissions,omitempty" gorm:"foreignKey:MemberID"`
	FinancialSettings *MemberFinancialSettings `json:"financial_settings,omitempty" gorm:"foreignKey:MemberID"`
	Relationships     *MemberRelationships     `json:"relationships,omitempty" gorm:"foreignKey:MemberID"`
	LoginLogs         []MemberLoginLogs        `json:"login_logs,omitempty" gorm:"foreignKey:MemberID"`
}

// TableName 指定表名
func (InternalMember) TableName() string {
	return "internal_members"
}

// MemberPermissions 内部成员权限表
type MemberPermissions struct {
	PermissionID   uint      `json:"permission_id" gorm:"primaryKey;column:permission_id"`
	MemberID       uint      `json:"member_id" gorm:"uniqueIndex;not null;comment:成员ID"`
	IsAuditor      bool      `json:"is_auditor" gorm:"default:false;comment:是否可审核"`
	CanReport      bool      `json:"can_report" gorm:"default:true;comment:是否可报单"`
	CanAcceptOrder bool      `json:"can_accept_order" gorm:"default:true;comment:是否可接单"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联关系
	Member *InternalMember `json:"member,omitempty" gorm:"foreignKey:MemberID"`
}

// TableName 指定表名
func (MemberPermissions) TableName() string {
	return "member_permissions"
}

// MemberFinancialSettings 内部成员财务设置表
type MemberFinancialSettings struct {
	SettingID      uint      `json:"setting_id" gorm:"primaryKey;column:setting_id"`
	MemberID       uint      `json:"member_id" gorm:"uniqueIndex;not null;comment:成员ID"`
	CommissionRate float64   `json:"commission_rate" gorm:"type:decimal(5,2);default:0.00;comment:提成比例（百分比）"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联关系
	Member *InternalMember `json:"member,omitempty" gorm:"foreignKey:MemberID"`
}

// TableName 指定表名
func (MemberFinancialSettings) TableName() string {
	return "member_financial_settings"
}

// MemberRelationships 内部成员关系表
type MemberRelationships struct {
	RelationshipID uint       `json:"relationship_id" gorm:"primaryKey;column:relationship_id"`
	MemberID       uint       `json:"member_id" gorm:"uniqueIndex;not null;comment:成员ID"`
	CreatorID      *uint      `json:"creator_id" gorm:"comment:录入人ID"`
	AssigneeID     *uint      `json:"assignee_id" gorm:"comment:归属人ID"`
	AssignedAt     *time.Time `json:"assigned_at" gorm:"comment:归属时间"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// 关联关系
	Member   *InternalMember `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Creator  *InternalMember `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	Assignee *InternalMember `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
}

// TableName 指定表名
func (MemberRelationships) TableName() string {
	return "member_relationships"
}

// MemberLoginLogs 内部成员登录日志表
type MemberLoginLogs struct {
	LogID       uint      `json:"log_id" gorm:"primaryKey;column:log_id"`
	MemberID    uint      `json:"member_id" gorm:"not null;comment:成员ID"`
	LoginTime   time.Time `json:"login_time" gorm:"comment:登录时间"`
	LoginIP     string    `json:"login_ip" gorm:"size:45;comment:登录IP"`
	UserAgent   string    `json:"user_agent" gorm:"type:text;comment:用户代理信息"`
	LoginStatus string    `json:"login_status" gorm:"type:enum('成功','失败');default:'成功';comment:登录状态"`

	// 关联关系
	Member *InternalMember `json:"member,omitempty" gorm:"foreignKey:MemberID"`
}

// TableName 指定表名
func (MemberLoginLogs) TableName() string {
	return "member_login_logs"
}

// 请求和响应结构
type LoginRequest struct {
	Account  string `json:"account" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}
type LoginResponse struct {
	Token     string          `json:"token"`
	ExpiresIn int64           `json:"expires_in"`
	User      *InternalMember `json:"user"`
}

type MemberCreateRequest struct {
	Account        string  `json:"account" binding:"required"`
	Password       string  `json:"password" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	PhoneNumber    string  `json:"phone_number"`
	Department     string  `json:"department"`
	UserRole       string  `json:"user_role" binding:"required"`
	Notes          string  `json:"notes"`
	IsAuditor      bool    `json:"is_auditor"`
	CanReport      bool    `json:"can_report"`
	CanAcceptOrder bool    `json:"can_accept_order"`  // 修复：与权限模型字段一致
	CommissionRate float64 `json:"commission_rate"`   // 修复：与财务设置模型字段一致
	CreatorID      *uint   `json:"creator_id"`
	AssigneeID     *uint   `json:"assignee_id"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	ExpiresIn    int64           `json:"expires_in"`
	User         *InternalMember `json:"user"`
}