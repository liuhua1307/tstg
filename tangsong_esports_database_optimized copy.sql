-- 唐宋电竞陪玩报单平台数据库设计 - 优化版本
-- 按业务粒度拆分复杂表结构
-- Created: 2025-01-11
-- Updated: 2025-01-11 - 使用自增主键

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =================================================================
-- 1. 内部成员相关表 (拆分 internal_members)
-- =================================================================

-- 1.1 内部成员基本信息表
DROP TABLE IF EXISTS `internal_members`;
CREATE TABLE `internal_members` (
  `member_id` INT NOT NULL AUTO_INCREMENT COMMENT '内部成员唯一ID',
  `account` VARCHAR(50) NOT NULL COMMENT '登录账号',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希值',
  `name` VARCHAR(50) NOT NULL COMMENT '姓名',
  `phone_number` VARCHAR(20) COMMENT '手机号码',
  `department` VARCHAR(50) COMMENT '所属部门',
  `user_role` VARCHAR(50) NOT NULL COMMENT '用户角色',
  `status` ENUM('正常', '禁用') NOT NULL DEFAULT '正常' COMMENT '账户状态',
  `is_enabled` BOOLEAN DEFAULT TRUE COMMENT '是否开启',
  `notes` TEXT COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`member_id`),
  UNIQUE KEY `uk_account` (`account`),
  UNIQUE KEY `uk_phone` (`phone_number`),
  KEY `idx_status` (`status`),
  KEY `idx_department` (`department`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内部成员基本信息表';

-- 1.2 内部成员权限表
DROP TABLE IF EXISTS `member_permissions`;
CREATE TABLE `member_permissions` (
  `permission_id` INT NOT NULL AUTO_INCREMENT COMMENT '权限记录ID',
  `member_id` INT NOT NULL COMMENT '成员ID',
  `is_auditor` BOOLEAN DEFAULT FALSE COMMENT '是否可审核',
  `can_report` BOOLEAN DEFAULT TRUE COMMENT '是否可报单',
  `can_accept_order` BOOLEAN DEFAULT TRUE COMMENT '是否可接单',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`permission_id`),
  UNIQUE KEY `uk_member_id` (`member_id`),
  CONSTRAINT `fk_member_permissions_member` FOREIGN KEY (`member_id`) REFERENCES `internal_members` (`member_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内部成员权限表';

-- 1.3 内部成员财务设置表
DROP TABLE IF EXISTS `member_financial_settings`;
CREATE TABLE `member_financial_settings` (
  `setting_id` INT NOT NULL AUTO_INCREMENT COMMENT '设置记录ID',
  `member_id` INT NOT NULL COMMENT '成员ID',
  `commission_rate` DECIMAL(5,2) DEFAULT 0.00 COMMENT '提成比例（百分比）',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`setting_id`),
  UNIQUE KEY `uk_member_id` (`member_id`),
  CONSTRAINT `fk_member_financial_member` FOREIGN KEY (`member_id`) REFERENCES `internal_members` (`member_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内部成员财务设置表';

-- 1.4 内部成员关系表
DROP TABLE IF EXISTS `member_relationships`;
CREATE TABLE `member_relationships` (
  `relationship_id` INT NOT NULL AUTO_INCREMENT COMMENT '关系记录ID',
  `member_id` INT NOT NULL COMMENT '成员ID',
  `creator_id` INT COMMENT '录入人ID',
  `assignee_id` INT COMMENT '归属人ID',
  `assigned_at` DATETIME COMMENT '归属时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`relationship_id`),
  UNIQUE KEY `uk_member_id` (`member_id`),
  KEY `idx_creator` (`creator_id`),
  KEY `idx_assignee` (`assignee_id`),
  CONSTRAINT `fk_member_relationships_member` FOREIGN KEY (`member_id`) REFERENCES `internal_members` (`member_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_member_relationships_creator` FOREIGN KEY (`creator_id`) REFERENCES `internal_members` (`member_id`),
  CONSTRAINT `fk_member_relationships_assignee` FOREIGN KEY (`assignee_id`) REFERENCES `internal_members` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内部成员关系表';

-- 1.5 内部成员登录记录表
DROP TABLE IF EXISTS `member_login_logs`;
CREATE TABLE `member_login_logs` (
  `log_id` INT NOT NULL AUTO_INCREMENT COMMENT '日志记录ID',
  `member_id` INT NOT NULL COMMENT '成员ID',
  `login_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `login_ip` VARCHAR(45) COMMENT '登录IP',
  `user_agent` TEXT COMMENT '用户代理信息',
  `login_status` ENUM('成功', '失败') DEFAULT '成功' COMMENT '登录状态',
  PRIMARY KEY (`log_id`),
  KEY `idx_member_time` (`member_id`, `login_time`),
  KEY `idx_login_time` (`login_time`),
  CONSTRAINT `fk_member_login_logs_member` FOREIGN KEY (`member_id`) REFERENCES `internal_members` (`member_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内部成员登录记录表';

-- =================================================================
-- 2. 客户相关表 (拆分 customers)
-- =================================================================

-- 2.1 客户基本信息表
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers` (
  `customer_id` INT NOT NULL AUTO_INCREMENT COMMENT '客户唯一ID',
  `account` VARCHAR(50) NOT NULL COMMENT '登录账号',
  `customer_name` VARCHAR(100) NOT NULL COMMENT '客户名称/昵称',
  `contact_method` VARCHAR(100) COMMENT '联系方式（微信/QQ/手机）',
  `phone_number` VARCHAR(20) COMMENT '手机号码',
  `member_birthday` DATE COMMENT '会员生日',
  `room_code` VARCHAR(50) COMMENT '房间码',
  `additional_info1` TEXT COMMENT '附加信息1',
  `additional_info2` TEXT COMMENT '附加信息2',
  `additional_info3` TEXT COMMENT '附加信息3',
  `notes` TEXT COMMENT '备注',
  `status` ENUM('正常', '禁用', '过期') NOT NULL DEFAULT '正常' COMMENT '客户账户状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '录入时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`customer_id`),
  UNIQUE KEY `uk_account` (`account`),
  UNIQUE KEY `uk_phone` (`phone_number`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户基本信息表';

-- 2.2 客户财务信息表
DROP TABLE IF EXISTS `customer_financial_info`;
CREATE TABLE `customer_financial_info` (
  `financial_id` INT NOT NULL AUTO_INCREMENT COMMENT '财务记录ID',
  `customer_id` INT NOT NULL COMMENT '客户ID',
  `initial_real_charge` DECIMAL(10,2) DEFAULT 0.00 COMMENT '初始实充金额',
  `total_consumption` DECIMAL(10,2) DEFAULT 0.00 COMMENT '历史消费总额',
  `total_real_charge` DECIMAL(10,2) DEFAULT 0.00 COMMENT '历史实际充值总额',
  `current_balance` DECIMAL(10,2) DEFAULT 0.00 COMMENT '当前可用余额',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`financial_id`),
  UNIQUE KEY `uk_customer_id` (`customer_id`),
  KEY `idx_balance` (`current_balance`),
  CONSTRAINT `fk_customer_financial_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户财务信息表';

-- 2.3 客户服务偏好表
DROP TABLE IF EXISTS `customer_preferences`;
CREATE TABLE `customer_preferences` (
  `preference_id` INT NOT NULL AUTO_INCREMENT COMMENT '偏好记录ID',
  `customer_id` INT NOT NULL COMMENT '客户ID',
  `exclusive_discount_type` ENUM('无折扣', '固定折扣') DEFAULT '无折扣' COMMENT '专属折扣类型',
  `platform_boss` VARCHAR(100) COMMENT '所属平台老板',
  `exclusive_cs` VARCHAR(100) COMMENT '专属服务客服',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`preference_id`),
  UNIQUE KEY `uk_customer_id` (`customer_id`),
  KEY `idx_platform_boss` (`platform_boss`),
  CONSTRAINT `fk_customer_preferences_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户服务偏好表';

-- =================================================================
-- 3. 充值记录表 (保持不变，结构合理)
-- =================================================================

DROP TABLE IF EXISTS `customer_recharge_history`;
CREATE TABLE `customer_recharge_history` (
  `recharge_id` INT NOT NULL AUTO_INCREMENT COMMENT '充值记录唯一ID',
  `customer_id` INT NOT NULL COMMENT '客户ID',
  `real_charge_amount` DECIMAL(10,2) NOT NULL COMMENT '实充金额',
  `gift_amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '赠送金额',
  `total_recharge_amount` DECIMAL(10,2) NOT NULL COMMENT '本次充值总额',
  `payment_method` ENUM('微信', '支付宝', '银行转账', '平台', '内部', '其他') NOT NULL COMMENT '付款方式',
  `transaction_id` TEXT COMMENT '收款单号/交易ID',
  `recharge_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '充值时间',
  `notes` TEXT COMMENT '备注',
  `operator_id` INT NOT NULL COMMENT '操作员ID',
  PRIMARY KEY (`recharge_id`),
  KEY `idx_customer_time` (`customer_id`, `recharge_at`),
  KEY `idx_operator` (`operator_id`),
  KEY `idx_payment_method` (`payment_method`),
  CONSTRAINT `fk_recharge_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
  CONSTRAINT `fk_recharge_operator` FOREIGN KEY (`operator_id`) REFERENCES `internal_members` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户充值记录表';

-- =================================================================
-- 4. 订单类别表 (保持不变，结构合理)
-- =================================================================

DROP TABLE IF EXISTS `order_categories`;
CREATE TABLE `order_categories` (
  `category_id` INT NOT NULL AUTO_INCREMENT COMMENT '订单类别唯一ID',
  `category_name` VARCHAR(100) NOT NULL COMMENT '订单类别名称',
  `sort_order` INT DEFAULT 0 COMMENT '排序序号',
  `is_active` BOOLEAN DEFAULT TRUE COMMENT '是否开启',
  `usage_scenario` VARCHAR(255) COMMENT '使用场景',
  `commission_rate` DECIMAL(5,2) DEFAULT 0.00 COMMENT '抽成比例（百分比）',
  `is_participating` BOOLEAN DEFAULT TRUE COMMENT '是否参与',
  `is_required` BOOLEAN DEFAULT FALSE COMMENT '是否必填',
  `is_accelerated` BOOLEAN DEFAULT FALSE COMMENT '是否加速',
  `additional_info` TEXT COMMENT '附加信息',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`category_id`),
  UNIQUE KEY `uk_category_name` (`category_name`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单类别表';

-- =================================================================
-- 5. 陪玩报单相关表 (拆分 playmate_orders)
-- =================================================================

-- 5.1 陪玩报单基本信息表
DROP TABLE IF EXISTS `playmate_orders`;
CREATE TABLE `playmate_orders` (
  `order_id` INT NOT NULL AUTO_INCREMENT COMMENT '订单唯一ID',
  `reporter_id` INT NOT NULL COMMENT '报单人ID',
  `customer_id` INT NOT NULL COMMENT '客户ID',
  `order_category_id` INT NOT NULL COMMENT '订单类别ID',
  `game` VARCHAR(100) NOT NULL COMMENT '游戏名称',
  `project_category` VARCHAR(100) NOT NULL COMMENT '项目分类',
  `playmate_level` VARCHAR(50) COMMENT '陪玩等级',
  `start_time` DATETIME NOT NULL COMMENT '开始时间',
  `end_time` DATETIME NOT NULL COMMENT '结束时间',
  `duration_hours` DECIMAL(5,2) NOT NULL COMMENT '陪玩时长（小时）',
  `is_teammate` BOOLEAN DEFAULT FALSE COMMENT '是否队友',
  `mode` VARCHAR(50) COMMENT '报单模式',
  `service_additional_info` TEXT COMMENT '服务附加说明',
  `internal_notes` TEXT COMMENT '内部备注',
  `order_notes` TEXT COMMENT '订单备注',
  `platform_owner` VARCHAR(100) COMMENT '所属平台老板',
  `report_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '报单时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`order_id`),
  KEY `idx_reporter` (`reporter_id`),
  KEY `idx_customer` (`customer_id`),
  KEY `idx_category` (`order_category_id`),
  KEY `idx_game` (`game`),
  KEY `idx_report_time` (`report_time`),
  CONSTRAINT `fk_playmate_orders_reporter` FOREIGN KEY (`reporter_id`) REFERENCES `internal_members` (`member_id`),
  CONSTRAINT `fk_playmate_orders_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
  CONSTRAINT `fk_playmate_orders_category` FOREIGN KEY (`order_category_id`) REFERENCES `order_categories` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='陪玩报单基本信息表';

-- 5.2 订单价格信息表
DROP TABLE IF EXISTS `order_pricing`;
CREATE TABLE `order_pricing` (
  `pricing_id` INT NOT NULL AUTO_INCREMENT COMMENT '价格记录ID',
  `order_id` INT NOT NULL COMMENT '订单ID',
  `unit_price` DECIMAL(10,2) NOT NULL COMMENT '单价（元/小时）',
  `total_price` DECIMAL(10,2) NOT NULL COMMENT '订单总价',
  `discount_amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '折扣总价',
  `final_price` DECIMAL(10,2) NOT NULL COMMENT '最终结算价格',
  `exclusive_discount` BOOLEAN DEFAULT FALSE COMMENT '是否专属折扣',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`pricing_id`),
  UNIQUE KEY `uk_order_id` (`order_id`),
  CONSTRAINT `fk_order_pricing_order` FOREIGN KEY (`order_id`) REFERENCES `playmate_orders` (`order_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单价格信息表';

-- 5.3 订单工作流状态表
DROP TABLE IF EXISTS `order_workflow`;
CREATE TABLE `order_workflow` (
  `workflow_id` INT NOT NULL AUTO_INCREMENT COMMENT '工作流记录ID',
  `order_id` INT NOT NULL COMMENT '订单ID',
  `order_status` ENUM('待处理', '驳回', '已确认', '已结算', '已完成', '已退回') NOT NULL DEFAULT '待处理' COMMENT '订单状态',
  `settlement_status` ENUM('未结算', '已结算') DEFAULT '未结算' COMMENT '结算状态',
  `approver_id` INT COMMENT '审批人ID',
  `approval_time` DATETIME COMMENT '审批时间',
  `rejection_reason` TEXT COMMENT '驳回原因',
  `settlement_time` DATETIME COMMENT '结算时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`workflow_id`),
  UNIQUE KEY `uk_order_id` (`order_id`),
  KEY `idx_order_status` (`order_status`),
  KEY `idx_settlement_status` (`settlement_status`),
  KEY `idx_approver` (`approver_id`),
  CONSTRAINT `fk_order_workflow_order` FOREIGN KEY (`order_id`) REFERENCES `playmate_orders` (`order_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_order_workflow_approver` FOREIGN KEY (`approver_id`) REFERENCES `internal_members` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单工作流状态表';

-- 5.4 订单支付信息表
DROP TABLE IF EXISTS `order_payment_info`;
CREATE TABLE `order_payment_info` (
  `payment_info_id` INT NOT NULL AUTO_INCREMENT COMMENT '支付信息记录ID',
  `order_id` INT NOT NULL COMMENT '订单ID',
  `payment_transaction_id` TEXT COMMENT '付款流水号',
  `payment_method` VARCHAR(50) COMMENT '付款方式',
  `payment_time` DATETIME COMMENT '付款时间',
  `payment_amount` DECIMAL(10,2) COMMENT '付款金额',
  `payment_status` ENUM('待付款', '已付款', '付款失败', '已退款') DEFAULT '待付款' COMMENT '付款状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`payment_info_id`),
  UNIQUE KEY `uk_order_id` (`order_id`),
  KEY `idx_payment_status` (`payment_status`),
  KEY `idx_payment_time` (`payment_time`),
  CONSTRAINT `fk_order_payment_order` FOREIGN KEY (`order_id`) REFERENCES `playmate_orders` (`order_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单支付信息表';

-- =================================================================
-- 6. 订单图片表 (保持不变，结构合理)
-- =================================================================

DROP TABLE IF EXISTS `order_images`;
CREATE TABLE `order_images` (
  `image_id` INT NOT NULL AUTO_INCREMENT COMMENT '图片唯一ID',
  `order_id` INT NOT NULL COMMENT '关联的订单ID',
  `image_url` VARCHAR(255) NOT NULL COMMENT '图片存储URL',
  `image_type` ENUM('结算截图', '聊天记录', '游戏截图', '其他') DEFAULT '其他' COMMENT '图片类型',
  `upload_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `notes` TEXT COMMENT '图片备注',
  PRIMARY KEY (`image_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_image_type` (`image_type`),
  CONSTRAINT `fk_order_images_order` FOREIGN KEY (`order_id`) REFERENCES `playmate_orders` (`order_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单图片表';

-- =================================================================
-- 7. 系统辅助表
-- =================================================================

-- 7.1 系统配置表
DROP TABLE IF EXISTS `system_configs`;
CREATE TABLE `system_configs` (
  `config_id` INT NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_key` VARCHAR(100) NOT NULL COMMENT '配置键',
  `config_value` TEXT COMMENT '配置值',
  `config_description` VARCHAR(255) COMMENT '配置描述',
  `is_active` BOOLEAN DEFAULT TRUE COMMENT '是否启用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`config_id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 7.2 操作日志表
DROP TABLE IF EXISTS `operation_logs`;
CREATE TABLE `operation_logs` (
  `log_id` INT NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `operator_id` INT NOT NULL COMMENT '操作人ID',
  `operation_type` VARCHAR(50) NOT NULL COMMENT '操作类型',
  `operation_module` VARCHAR(50) COMMENT '操作模块',
  `operation_description` TEXT COMMENT '操作描述',
  `target_id` VARCHAR(36) COMMENT '操作目标ID',
  `target_type` VARCHAR(50) COMMENT '操作目标类型',
  `old_data` JSON COMMENT '操作前数据',
  `new_data` JSON COMMENT '操作后数据',
  `ip_address` VARCHAR(45) COMMENT 'IP地址',
  `user_agent` TEXT COMMENT '用户代理',
  `operation_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`log_id`),
  KEY `idx_operator` (`operator_id`),
  KEY `idx_operation_type` (`operation_type`),
  KEY `idx_operation_time` (`operation_time`),
  KEY `idx_target` (`target_id`, `target_type`),
  CONSTRAINT `fk_operation_logs_operator` FOREIGN KEY (`operator_id`) REFERENCES `internal_members` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- =================================================================
-- 8. 初始化数据
-- =================================================================

-- 插入默认系统配置
INSERT INTO `system_configs` (`config_key`, `config_value`, `config_description`) VALUES
('default_commission_rate', '0.15', '默认提成比例'),
('max_order_hours', '24', '单次订单最大时长（小时）'),
('auto_settlement_enabled', 'false', '是否启用自动结算'),
('order_image_max_size', '5242880', '订单图片最大尺寸（字节）'),
('platform_name', '唐宋电竞陪玩平台', '平台名称');

-- 插入默认订单类别
INSERT INTO `order_categories` (`category_name`, `sort_order`, `commission_rate`, `usage_scenario`) VALUES
('英雄联盟', 10, 0.20, '排位赛、匹配、大乱斗'),
('三角洲行动', 20, 0.18, 'PVP、PVE模式'),
('瓦罗兰特', 30, 0.22, '竞技模式'),
('明星陪', 40, 0.25, '高端陪玩服务'),
('拍卖单', 50, 0.15, '特殊拍卖服务');

SET FOREIGN_KEY_CHECKS = 1;

-- =================================================================
-- 结构优化说明
-- =================================================================

/*
主要优化内容：

1. 内部成员表拆分：
   - internal_members: 基本信息
   - member_permissions: 权限管理
   - member_financial_settings: 财务设置
   - member_relationships: 成员关系
   - member_login_logs: 登录记录

2. 客户表拆分：
   - customers: 基本信息
   - customer_financial_info: 财务信息
   - customer_preferences: 服务偏好

3. 陪玩报单表拆分：
   - playmate_orders: 基本信息
   - order_pricing: 价格信息
   - order_workflow: 工作流状态
   - order_payment_info: 支付信息

4. 新增辅助表：
   - system_configs: 系统配置
   - operation_logs: 操作日志

5. 主键优化：
   - 所有表主键改为INT AUTO_INCREMENT
   - 提高性能，减少存储空间
   - 简化外键关系

优势：
- 表结构更清晰，职责单一
- 提高查询性能
- 便于扩展和维护
- 减少数据冗余
- 支持更细粒度的权限控制
- 自增主键提供更好的性能
*/ 