# 唐宋电竞陪玩报单平台 API 测试指南

## 基础信息
- **基础URL**: `http://localhost:8080/api/v1`
- **默认管理员**: `admin` / `123456`

## 1. 用户认证

### 登录
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "account": "admin",
    "password": "123456"
  }'
```

**响应示例:**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user": {
      "member_id": 1,
      "account": "admin",
      "name": "系统管理员"
    }
  }
}
```

## 2. 内部成员管理

### 获取成员列表
```bash
curl -X GET "http://localhost:8080/api/v1/members?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 创建成员
```bash
curl -X POST http://localhost:8080/api/v1/members \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test001",
    "password": "123456",
    "name": "测试陪玩",
    "phone_number": "13800138000",
    "department": "陪玩部",
    "user_role": "陪玩",
    "notes": "测试账号",
    "is_auditor": false,
    "can_report": true,
    "can_accept_order": true,
    "commission_rate": 0.15
  }'
```

### 获取成员详情
```bash
curl -X GET http://localhost:8080/api/v1/members/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 3. 客户管理

### 获取客户列表
```bash
curl -X GET "http://localhost:8080/api/v1/customers?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 创建客户
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "customer001",
    "customer_name": "测试客户",
    "contact_method": "微信: test_wx",
    "phone_number": "13900139000",
    "member_birthday": "1990-01-01",
    "room_code": "ROOM001",
    "notes": "VIP客户",
    "initial_real_charge": 1000.00,
    "exclusive_discount_type": "固定折扣",
    "platform_boss": "平台老板A",
    "exclusive_cs": "客服小王"
  }'
```

### 客户充值
```bash
curl -X POST http://localhost:8080/api/v1/customers/1/recharge \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "real_charge_amount": 500.00,
    "gift_amount": 50.00,
    "payment_method": "微信",
    "transaction_id": "WX20250715001",
    "notes": "充值测试"
  }'
```

## 4. 订单类别管理

### 获取订单类别
```bash
curl -X GET http://localhost:8080/api/v1/order-categories \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 创建订单类别
```bash
curl -X POST http://localhost:8080/api/v1/order-categories \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_name": "王者荣耀",
    "sort_order": 60,
    "usage_scenario": "排位、巅峰赛",
    "commission_rate": 0.20,
    "is_participating": true,
    "is_required": false,
    "is_accelerated": false,
    "additional_info": "热门游戏"
  }'
```

## 5. 陪玩订单管理

### 获取订单列表
```bash
curl -X GET "http://localhost:8080/api/v1/orders?page=1&page_size=10&order_status=待处理" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 创建订单
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "order_category_id": 1,
    "game": "英雄联盟",
    "project_category": "排位赛",
    "playmate_level": "黄金",
    "start_time": "2025-07-15 20:00:00",
    "end_time": "2025-07-15 22:00:00",
    "duration_hours": 2.0,
    "unit_price": 50.00,
    "is_teammate": false,
    "mode": "单排",
    "service_additional_info": "要求有语音",
    "internal_notes": "客户比较挑剔",
    "order_notes": "上分专用",
    "platform_owner": "平台老板A",
    "exclusive_discount": false
  }'
```

### 更新订单状态
```bash
curl -X PUT http://localhost:8080/api/v1/orders/1/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_status": "已确认"
  }'
```

### 驳回订单
```bash
curl -X PUT http://localhost:8080/api/v1/orders/1/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_status": "驳回",
    "rejection_reason": "信息不完整，请补充客户需求"
  }'
```

## 6. 系统管理

### 获取系统配置
```bash
curl -X GET http://localhost:8080/api/v1/configs \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 更新配置
```bash
curl -X PUT http://localhost:8080/api/v1/configs/default_commission_rate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "config_value": "0.18",
    "config_description": "默认提成比例(已调整)"
  }'
```

### 获取操作日志
```bash
curl -X GET "http://localhost:8080/api/v1/logs?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 7. 测试流程

### 完整测试流程
1. **登录获取Token**
2. **创建陪玩成员**
3. **创建客户**
4. **为客户充值**
5. **创建订单**
6. **审批订单**
7. **查看操作日志**

### 状态流转测试
```bash
# 1. 创建订单 -> 待处理
# 2. 确认订单
curl -X PUT http://localhost:8080/api/v1/orders/1/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_status": "已确认"}'

# 3. 结算订单
curl -X PUT http://localhost:8080/api/v1/orders/1/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_status": "已结算"}'
```

## 8. 错误处理

### 常见错误码
- `200`: 成功
- `401`: 未授权（Token无效或过期）
- `403`: 禁止访问
- `404`: 资源不存在
- `500`: 服务器内部错误

### 错误响应格式
```json
{
  "code": 500,
  "message": "错误描述"
}
```

## 9. 数据格式说明

### 时间格式
- 日期时间: `2025-07-15 20:00:00`
- 日期: `2025-07-15`

### 订单状态
- `待处理`: 新创建的订单
- `已确认`: 审核通过的订单
- `驳回`: 审核未通过的订单
- `已结算`: 已完成结算的订单
- `已完成`: 订单全部完成
- `已退回`: 订单被退回

### 支付方式
- `微信`
- `支付宝`
- `银行转账`
- `平台`
- `内部`
- `其他` 