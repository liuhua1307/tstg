package database

import (
	"fmt"
	"log"
	"tangsong-esports/config"
	"tangsong-esports/models"
	"tangsong-esports/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})

	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	log.Println("数据库连接成功")

	// 按依赖顺序自动迁移数据表
	log.Println("开始数据库迁移...")

	// 第一步：创建基础表（不依赖其他表的表）
	err = DB.AutoMigrate(
		&models.InternalMember{},
		&models.Customer{},
		&models.OrderCategory{},
		&models.SystemConfig{},
	)
	if err != nil {
		log.Fatal("基础表迁移失败:", err)
	}

	// 第二步：创建依赖表
	err = DB.AutoMigrate(
		&models.MemberPermissions{},
		&models.MemberFinancialSettings{},
		&models.MemberRelationships{},
		&models.MemberLoginLogs{},
		&models.CustomerFinancialInfo{},
		&models.CustomerPreferences{},
		&models.CustomerRechargeHistory{},
		&models.PlaymateOrder{},
		&models.OperationLog{},
	)
	if err != nil {
		log.Fatal("依赖表迁移失败:", err)
	}

	// 第三步：创建订单相关表
	err = DB.AutoMigrate(
		&models.OrderPricing{},
		&models.OrderWorkflow{},
		&models.OrderPaymentInfo{},
		&models.OrderImages{},
	)
	if err != nil {
		log.Fatal("订单表迁移失败:", err)
	}

	log.Println("数据库迁移完成")

	// 初始化基础数据
	initBaseData()
}

func GetDB() *gorm.DB {
	return DB
}

// 初始化基础数据
func initBaseData() {
	// 检查表是否存在，如果不存在则跳过
	if !DB.Migrator().HasTable(&models.SystemConfig{}) {
		log.Println("数据表不存在，跳过基础数据初始化")
		return
	}

	// 插入默认系统配置
	configs := []models.SystemConfig{
		{ConfigKey: "default_commission_rate", ConfigValue: "0.15", ConfigDescription: "默认提成比例"},
		{ConfigKey: "max_order_hours", ConfigValue: "24", ConfigDescription: "单次订单最大时长（小时）"},
		{ConfigKey: "auto_settlement_enabled", ConfigValue: "false", ConfigDescription: "是否启用自动结算"},
		{ConfigKey: "order_image_max_size", ConfigValue: "5242880", ConfigDescription: "订单图片最大尺寸（字节）"},
		{ConfigKey: "platform_name", ConfigValue: "唐宋电竞陪玩平台", ConfigDescription: "平台名称"},
	}

	for _, config := range configs {
		var count int64
		DB.Model(&models.SystemConfig{}).Where("config_key = ?", config.ConfigKey).Count(&count)
		if count == 0 {
			DB.Create(&config)
		}
	}

	// 插入默认订单类别
	if DB.Migrator().HasTable(&models.OrderCategory{}) {
		categories := []models.OrderCategory{
			{CategoryName: "英雄联盟", SortOrder: 10, CommissionRate: 0.20, UsageScenario: "排位赛、匹配、大乱斗"},
			{CategoryName: "三角洲行动", SortOrder: 20, CommissionRate: 0.18, UsageScenario: "PVP、PVE模式"},
			{CategoryName: "瓦罗兰特", SortOrder: 30, CommissionRate: 0.22, UsageScenario: "竞技模式"},
			{CategoryName: "明星陪", SortOrder: 40, CommissionRate: 0.25, UsageScenario: "高端陪玩服务"},
			{CategoryName: "拍卖单", SortOrder: 50, CommissionRate: 0.15, UsageScenario: "特殊拍卖服务"},
		}

		for _, category := range categories {
			var count int64
			DB.Model(&models.OrderCategory{}).Where("category_name = ?", category.CategoryName).Count(&count)
			if count == 0 {
				DB.Create(&category)
			}
		}
	}

	// 创建默认管理员账户（如果表存在且不存在管理员）
	if DB.Migrator().HasTable(&models.InternalMember{}) {
		var count int64
		DB.Model(&models.InternalMember{}).Where("account = ?", "admin").Count(&count)
		if count == 0 {
			// 生成密码哈希
			hashedPassword, err := utils.HashPassword("123456")
			if err != nil {
				log.Printf("密码哈希生成失败: %v", err)
				hashedPassword = "123456" // 回退到明文（仅用于开发）
			}

			// 创建默认管理员
			admin := models.InternalMember{
				Account:      "admin",
				PasswordHash: hashedPassword,
				Name:         "系统管理员",
				UserRole:     "超级管理员",
				Status:       "正常",
				IsEnabled:    true,
				Notes:        "默认管理员账户",
			}
			DB.Create(&admin)
			log.Println("创建默认管理员账户: admin / 123456")
		}
	}

	log.Println("基础数据初始化完成")
}
