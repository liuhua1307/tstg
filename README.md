# 唐宋电竞陪玩报单平台后端

基于 Gin 框架开发的唐宋电竞陪玩报单平台后端API系统。

## 功能特性

- 🔐 **用户认证** - JWT认证系统
- 👥 **内部成员管理** - 支持多角色权限管理
- 👤 **客户管理** - 客户信息、余额、充值管理
- 📋 **订单管理** - 陪玩报单、审批、状态管理
- 📊 **财务管理** - 提成计算、结算记录
- 🏷️ **订单分类** - 灵活的订单类别配置
- 📸 **文件上传** - 订单图片上传
- 📜 **操作日志** - 完整的操作审计
- ⚙️ **系统配置** - 可配置的系统参数

## 技术栈

- **Web框架**: Gin
- **数据库**: MySQL + GORM
- **认证**: JWT
- **配置管理**: Viper
- **API文档**: Swagger
- **密码加密**: bcrypt

## 项目结构

```
.
├── main.go                 # 程序入口
├── config/                 # 配置管理
│   └── config.go
├── models/                 # 数据模型
│   ├── member.go          # 内部成员相关模型
│   ├── customer.go        # 客户相关模型
│   ├── order.go           # 订单相关模型
│   └── system.go          # 系统相关模型
├── database/               # 数据库相关
│   └── database.go
├── controllers/            # 控制器
│   ├── auth.go            # 认证控制器
│   ├── member.go          # 成员管理
│   ├── customer.go        # 客户管理
│   ├── order.go           # 订单管理
│   └── system.go          # 系统管理
├── middleware/             # 中间件
│   └── auth.go            # 认证中间件
├── router/                 # 路由配置
│   └── router.go
├── utils/                  # 工具类
│   ├── jwt.go             # JWT工具
│   ├── password.go        # 密码工具
│   ├── response.go        # 响应工具
│   └── logger.go          # 日志工具
├── config.yaml            # 配置文件
└── go.mod                 # 依赖管理
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

修改 `config.yaml` 文件中的数据库配置：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: tangsong_esports
  charset: utf8mb4
```

### 4. 创建数据库

```sql
CREATE DATABASE tangsong_esports CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 运行项目

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

### 6. 查看API文档

访问 `http://localhost:8080/swagger/index.html` 查看Swagger API文档。

## API 端点

### 认证
- `POST /api/v1/login` - 用户登录

### 内部成员管理
- `GET /api/v1/members` - 获取成员列表
- `POST /api/v1/members` - 创建成员
- `PUT /api/v1/members/:id` - 更新成员
- `DELETE /api/v1/members/:id` - 删除成员
- `GET /api/v1/members/:id` - 获取成员详情

### 客户管理
- `GET /api/v1/customers` - 获取客户列表
- `POST /api/v1/customers` - 创建客户
- `PUT /api/v1/customers/:id` - 更新客户
- `DELETE /api/v1/customers/:id` - 删除客户
- `GET /api/v1/customers/:id` - 获取客户详情
- `POST /api/v1/customers/:id/recharge` - 客户充值

### 订单管理
- `GET /api/v1/orders` - 获取订单列表
- `POST /api/v1/orders` - 创建订单
- `PUT /api/v1/orders/:id` - 更新订单
- `GET /api/v1/orders/:id` - 获取订单详情
- `PUT /api/v1/orders/:id/status` - 更新订单状态

### 系统管理
- `GET /api/v1/configs` - 获取系统配置
- `PUT /api/v1/configs/:key` - 更新配置
- `GET /api/v1/logs` - 获取操作日志

## 数据库设计

项目采用优化的数据库设计，将业务数据按职责拆分到不同表中：

### 内部成员相关表
- `internal_members` - 成员基本信息
- `member_permissions` - 成员权限
- `member_financial_settings` - 财务设置
- `member_relationships` - 成员关系
- `member_login_logs` - 登录记录

### 客户相关表
- `customers` - 客户基本信息
- `customer_financial_info` - 客户财务信息
- `customer_preferences` - 客户偏好设置
- `customer_recharge_history` - 充值记录

### 订单相关表
- `order_categories` - 订单类别
- `playmate_orders` - 陪玩订单基本信息
- `order_pricing` - 订单价格信息
- `order_workflow` - 订单工作流状态
- `order_payment_info` - 订单支付信息
- `order_images` - 订单图片

### 系统相关表
- `system_configs` - 系统配置
- `operation_logs` - 操作日志

## 开发说明

### 添加新的API接口

1. 在 `models/` 中定义数据模型
2. 在 `controllers/` 中实现业务逻辑
3. 在 `router/router.go` 中添加路由
4. 添加必要的中间件和权限控制

### 配置管理

系统配置通过 `config.yaml` 文件管理，支持以下配置项：

- 服务器端口
- 数据库连接
- JWT密钥和过期时间
- 日志级别

### 权限控制

通过JWT认证和角色权限系统实现：

- 不同角色拥有不同的操作权限
- 中间件自动验证用户身份
- 支持细粒度的权限控制

## 部署

### 生产环境配置

1. 修改 `config.yaml` 中的 `mode` 为 `release`
2. 配置生产数据库连接
3. 设置安全的JWT密钥
4. 配置适当的日志级别

### Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
CMD ["./main"]
```

## 许可证

MIT License 