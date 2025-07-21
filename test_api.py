#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
唐宋电竞陪玩报单平台 API 自动化测试脚本
运行方式: python test_api.py
"""

import requests
import json
import time
import sys
from typing import Dict, Any, Optional
from dataclasses import dataclass
from datetime import datetime, timedelta


@dataclass
class TestResult:
    """测试结果数据类"""
    endpoint: str
    method: str
    status_code: int
    success: bool
    response_time: float
    error_message: Optional[str] = None
    response_data: Optional[Dict] = None


class APITester:
    """API测试类"""
    
    def __init__(self, base_url: str = "http://localhost:8080/api/v1"):
        self.base_url = base_url
        self.token = None
        self.session = requests.Session()
        self.test_results = []
        self.test_data = {}  # 存储测试过程中创建的数据ID
        
        # 设置请求头
        self.session.headers.update({
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        })
    
    def log(self, message: str, level: str = "INFO"):
        """日志输出"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] [{level}] {message}")
    
    def make_request(self, method: str, endpoint: str, data: Dict = None, 
                    params: Dict = None, auth_required: bool = True) -> TestResult:
        """发送HTTP请求并记录结果"""
        url = f"{self.base_url}{endpoint}"
        
        # 设置认证头
        headers = {}
        if auth_required and self.token:
            headers['Authorization'] = f'Bearer {self.token}'
        
        start_time = time.time()
        
        try:
            if method.upper() == 'GET':
                response = self.session.get(url, params=params, headers=headers)
            elif method.upper() == 'POST':
                response = self.session.post(url, json=data, headers=headers)
            elif method.upper() == 'PUT':
                response = self.session.put(url, json=data, headers=headers)
            elif method.upper() == 'DELETE':
                response = self.session.delete(url, headers=headers)
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")
            
            response_time = time.time() - start_time
            
            # 尝试解析JSON响应
            try:
                response_data = response.json()
            except json.JSONDecodeError:
                response_data = {"text": response.text}
            
            success = 200 <= response.status_code < 300
            error_message = None if success else response_data.get('message', 'Unknown error')
            
            result = TestResult(
                endpoint=endpoint,
                method=method.upper(),
                status_code=response.status_code,
                success=success,
                response_time=response_time,
                error_message=error_message,
                response_data=response_data
            )
            
            self.test_results.append(result)
            
            status = "✅ PASS" if success else "❌ FAIL"
            self.log(f"{status} {method.upper()} {endpoint} ({response.status_code}) - {response_time:.3f}s")
            
            if not success:
                self.log(f"Error: {error_message}", "ERROR")
            
            return result
            
        except Exception as e:
            response_time = time.time() - start_time
            result = TestResult(
                endpoint=endpoint,
                method=method.upper(),
                status_code=0,
                success=False,
                response_time=response_time,
                error_message=str(e)
            )
            
            self.test_results.append(result)
            self.log(f"❌ FAIL {method.upper()} {endpoint} - Exception: {str(e)}", "ERROR")
            
            return result
    
    def login(self) -> bool:
        """用户登录获取Token"""
        self.log("开始登录测试...")
        
        login_data = {
            "account": "admin",
            "password": "123456"
        }
        
        result = self.make_request('POST', '/login', data=login_data, auth_required=False)
        
        if result.success and result.response_data:
            token_data = result.response_data.get('data', {})
            self.token = token_data.get('token')
            
            if self.token:
                self.log(f"登录成功，获取到Token: {self.token[:50]}...")
                return True
            else:
                self.log("登录响应中未找到token", "ERROR")
                return False
        else:
            self.log("登录失败", "ERROR")
            return False
    
    def test_members(self):
        """测试成员管理接口"""
        self.log("开始测试成员管理接口...")
        
        # 1. 获取成员列表
        self.make_request('GET', '/members', params={'page': 1, 'page_size': 10})
        
        # 2. 搜索成员
        self.make_request('GET', '/members', params={'account': 'admin'})
        
        # 3. 创建成员
        member_data = {
            "account": f"test_{int(time.time())}",
            "password": "123456",
            "name": "测试陪玩",
            "phone_number": "13800138000",
            "department": "陪玩部",
            "user_role": "陪玩",
            "notes": "自动化测试创建",
            "is_auditor": False,
            "can_report": True,
            "can_accept_order": True,
            "commission_rate": 0.15,
            "creator_id": 1
        }
        
        result = self.make_request('POST', '/members', data=member_data)
        if result.success and result.response_data:
            member_id = result.response_data.get('data', {}).get('member_id')
            if member_id:
                self.test_data['member_id'] = member_id
                
                # 4. 获取成员详情
                self.make_request('GET', f'/members/{member_id}')
                
                # 5. 更新成员
                updated_data = member_data.copy()
                updated_data['name'] = '更新后的测试陪玩'
                updated_data['password'] = ''  # 不更新密码
                self.make_request('PUT', f'/members/{member_id}', data=updated_data)
                
                # 6. 删除成员（软删除）
                self.make_request('DELETE', f'/members/{member_id}')
    
    def test_customers(self):
        """测试客户管理接口"""
        self.log("开始测试客户管理接口...")
        
        # 1. 获取客户列表
        self.make_request('GET', '/customers', params={'page': 1, 'page_size': 10})
        
        # 2. 创建客户
        customer_data = {
            "account": f"customer_{int(time.time())}",
            "customer_name": "测试客户",
            "contact_method": "微信: test_wx",
            "phone_number": "13900139000",
            "member_birthday": "1990-01-01",
            "room_code": "ROOM001",
            "notes": "自动化测试创建",
            "initial_real_charge": 1000.00,
            "exclusive_discount_type": "固定折扣",
            "platform_boss": "平台老板A",
            "exclusive_cs": "客服小王"
        }
        
        result = self.make_request('POST', '/customers', data=customer_data)
        if result.success and result.response_data:
            customer_id = result.response_data.get('data', {}).get('customer_id')
            if customer_id:
                self.test_data['customer_id'] = customer_id
                
                # 3. 获取客户详情
                self.make_request('GET', f'/customers/{customer_id}')
                
                # 4. 更新客户
                updated_data = customer_data.copy()
                updated_data['customer_name'] = '更新后的测试客户'
                self.make_request('PUT', f'/customers/{customer_id}', data=updated_data)
                
                # 5. 客户充值
                recharge_data = {
                    "real_charge_amount": 500.00,
                    "gift_amount": 50.00,
                    "payment_method": "微信",
                    "transaction_id": f"WX{int(time.time())}",
                    "notes": "自动化测试充值"
                }
                self.make_request('POST', f'/customers/{customer_id}/recharge', data=recharge_data)
                
                # 6. 获取充值记录
                self.make_request('GET', f'/customers/{customer_id}/recharge-history')
                
                # 7. 删除客户
                self.make_request('DELETE', f'/customers/{customer_id}')
    
    def test_order_categories(self):
        """测试订单类别管理接口"""
        self.log("开始测试订单类别管理接口...")
        
        # 1. 获取订单类别列表
        self.make_request('GET', '/order-categories')
        
        # 2. 创建订单类别
        category_data = {
            "category_name": f"测试游戏_{int(time.time())}",
            "sort_order": 100,
            "usage_scenario": "自动化测试",
            "commission_rate": 0.20,
            "is_participating": True,
            "is_required": False,
            "is_accelerated": False,
            "additional_info": "自动化测试创建"
        }
        
        result = self.make_request('POST', '/order-categories', data=category_data)
        if result.success and result.response_data:
            category_id = result.response_data.get('data', {}).get('category_id')
            if category_id:
                self.test_data['category_id'] = category_id
                
                # 3. 更新订单类别
                updated_data = category_data.copy()
                updated_data['category_name'] = f"更新后的测试游戏_{int(time.time())}"
                self.make_request('PUT', f'/order-categories/{category_id}', data=updated_data)
                
                # 4. 删除订单类别
                self.make_request('DELETE', f'/order-categories/{category_id}')
    
    def test_orders(self):
        """测试订单管理接口"""
        self.log("开始测试订单管理接口...")
        
        # 确保有必要的数据
        if 'customer_id' not in self.test_data or 'category_id' not in self.test_data:
            self.log("缺少必要的测试数据，跳过订单测试", "WARN")
            return
        
        # 1. 获取订单列表
        self.make_request('GET', '/orders', params={'page': 1, 'page_size': 10})
        
        # 2. 创建订单
        now = datetime.now()
        start_time = now.strftime("%Y-%m-%d %H:%M:%S")
        end_time = (now + timedelta(hours=2)).strftime("%Y-%m-%d %H:%M:%S")
        
        order_data = {
            "customer_id": self.test_data['customer_id'],
            "order_category_id": self.test_data['category_id'],
            "game": "测试游戏",
            "project_category": "自动化测试",
            "playmate_level": "测试",
            "start_time": start_time,
            "end_time": end_time,
            "duration_hours": 2.0,
            "unit_price": 50.00,
            "is_teammate": False,
            "mode": "测试模式",
            "service_additional_info": "自动化测试",
            "internal_notes": "自动化测试创建",
            "order_notes": "测试订单",
            "platform_owner": "测试平台",
            "exclusive_discount": False
        }
        
        result = self.make_request('POST', '/orders', data=order_data)
        if result.success and result.response_data:
            order_id = result.response_data.get('data', {}).get('order_id')
            if order_id:
                self.test_data['order_id'] = order_id
                
                # 3. 获取订单详情
                self.make_request('GET', f'/orders/{order_id}')
                
                # 4. 更新订单
                updated_data = order_data.copy()
                updated_data['game'] = '更新后的测试游戏'
                self.make_request('PUT', f'/orders/{order_id}', data=updated_data)
                
                # 5. 更新订单状态 - 确认
                status_data = {"order_status": "已确认"}
                self.make_request('PUT', f'/orders/{order_id}/status', data=status_data)
                
                # 6. 更新订单状态 - 驳回
                reject_data = {
                    "order_status": "驳回",
                    "rejection_reason": "自动化测试驳回"
                }
                self.make_request('PUT', f'/orders/{order_id}/status', data=reject_data)
        
        # 7. 按状态筛选订单
        self.make_request('GET', '/orders', params={'order_status': '待处理'})
        
        # 8. 按时间范围筛选订单
        today = datetime.now().strftime("%Y-%m-%d")
        self.make_request('GET', '/orders', params={
            'start_date': today,
            'end_date': today
        })
    
    def test_system_config(self):
        """测试系统配置接口"""
        self.log("开始测试系统配置接口...")
        
        # 1. 获取系统配置
        result = self.make_request('GET', '/configs')
        
        if result.success and result.response_data:
            configs = result.response_data.get('data', [])
            if configs:
                # 选择第一个配置进行更新测试
                config = configs[0]
                config_key = config.get('config_key')
                
                if config_key:
                    # 2. 更新配置
                    update_data = {
                        "config_value": "0.18",
                        "config_description": "自动化测试更新",
                        "is_active": True
                    }
                    self.make_request('PUT', f'/configs/{config_key}', data=update_data)
        
        # 3. 筛选启用的配置
        self.make_request('GET', '/configs', params={'is_active': True})
    
    def test_operation_logs(self):
        """测试操作日志接口"""
        self.log("开始测试操作日志接口...")
        
        # 1. 获取操作日志
        self.make_request('GET', '/logs', params={'page': 1, 'page_size': 10})
        
        # 2. 按操作类型筛选
        self.make_request('GET', '/logs', params={'operation_type': '创建'})
        
        # 3. 按模块筛选
        self.make_request('GET', '/logs', params={'operation_module': '订单管理'})
        
        # 4. 按时间范围筛选
        today = datetime.now().strftime("%Y-%m-%d")
        self.make_request('GET', '/logs', params={
            'start_date': today,
            'end_date': today
        })
    
    def run_all_tests(self):
        """运行所有测试"""
        self.log("开始API自动化测试...")
        self.log(f"测试目标: {self.base_url}")
        
        # 1. 登录测试
        if not self.login():
            self.log("登录失败，终止测试", "ERROR")
            return False
        
        # 2. 运行各模块测试
        try:
            self.test_members()
            self.test_customers()
            self.test_order_categories()
            self.test_orders()
            self.test_system_config()
            self.test_operation_logs()
        except Exception as e:
            self.log(f"测试过程中出现异常: {str(e)}", "ERROR")
            return False
        
        return True
    
    def generate_report(self):
        """生成测试报告"""
        self.log("生成测试报告...")
        
        total_tests = len(self.test_results)
        passed_tests = sum(1 for result in self.test_results if result.success)
        failed_tests = total_tests - passed_tests
        
        success_rate = (passed_tests / total_tests * 100) if total_tests > 0 else 0
        
        report = {
            "test_summary": {
                "total_tests": total_tests,
                "passed_tests": passed_tests,
                "failed_tests": failed_tests,
                "success_rate": f"{success_rate:.2f}%",
                "test_time": datetime.now().isoformat()
            },
            "test_results": []
        }
        
        # 按模块分组结果
        modules = {}
        for result in self.test_results:
            module = result.endpoint.split('/')[1] if '/' in result.endpoint else 'auth'
            if module not in modules:
                modules[module] = []
            modules[module].append({
                "endpoint": result.endpoint,
                "method": result.method,
                "status_code": result.status_code,
                "success": result.success,
                "response_time": f"{result.response_time:.3f}s",
                "error_message": result.error_message
            })
        
        report["test_results"] = modules
        
        # 保存报告到文件
        report_file = f"api_test_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        with open(report_file, 'w', encoding='utf-8') as f:
            json.dump(report, f, ensure_ascii=False, indent=2)
        
        # 打印摘要
        print("\n" + "="*60)
        print("API 测试报告摘要")
        print("="*60)
        print(f"总测试数: {total_tests}")
        print(f"通过测试: {passed_tests}")
        print(f"失败测试: {failed_tests}")
        print(f"成功率: {success_rate:.2f}%")
        print(f"报告文件: {report_file}")
        print("="*60)
        
        # 显示失败的测试
        if failed_tests > 0:
            print("\n失败的测试:")
            print("-" * 40)
            for result in self.test_results:
                if not result.success:
                    print(f"❌ {result.method} {result.endpoint} ({result.status_code})")
                    if result.error_message:
                        print(f"   错误: {result.error_message}")
        
        return report


def main():
    """主函数"""
    print("唐宋电竞陪玩报单平台 API 自动化测试")
    print("="*50)
    
    # 检查服务器是否运行
    tester = APITester()
    
    try:
        response = requests.get(f"{tester.base_url}/../swagger/index.html", timeout=5)
        if response.status_code != 200:
            print("⚠️  警告: 无法访问Swagger文档，请确认服务器是否正常运行")
    except requests.exceptions.RequestException:
        print("❌ 错误: 无法连接到服务器，请确认服务器是否在运行")
        print(f"   测试地址: {tester.base_url}")
        sys.exit(1)
    
    # 运行测试
    success = tester.run_all_tests()
    
    # 生成报告
    report = tester.generate_report()
    
    # 根据测试结果退出
    if success and report["test_summary"]["failed_tests"] == 0:
        print("\n🎉 所有测试通过!")
        sys.exit(0)
    else:
        print("\n❌ 部分测试失败，请检查报告详情")
        sys.exit(1)


if __name__ == "__main__":
    main() 