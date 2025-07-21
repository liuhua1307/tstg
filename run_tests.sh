#!/bin/bash

# 唐宋电竞陪玩报单平台 API 测试脚本
# 使用方法: ./run_tests.sh

echo "🚀 唐宋电竞陪玩报单平台 API 测试启动"
echo "=========================================="

# 检查Python是否安装
if ! command -v python3 &> /dev/null; then
    echo "❌ 错误: 未找到 python3，请先安装 Python 3.7+"
    exit 1
fi

# 检查并安装依赖
echo "📦 检查Python依赖..."
if [ -f "requirements.txt" ]; then
    pip3 install -r requirements.txt
else
    echo "⚠️  未找到 requirements.txt，尝试安装基础依赖..."
    pip3 install requests
fi

# 检查服务器是否运行
echo "🔍 检查服务器状态..."
if curl -s http://localhost:8080/api/v1/login > /dev/null; then
    echo "✅ 服务器正在运行"
else
    echo "❌ 错误: 无法连接到服务器"
    echo "   请确保服务器正在 http://localhost:8080 运行"
    echo "   启动服务器: go run main.go"
    exit 1
fi

# 运行API测试
echo "🧪 开始运行API测试..."
python3 test_api.py

# 检查测试结果
if [ $? -eq 0 ]; then
    echo "🎉 测试完成！所有测试通过"
else
    echo "❌ 测试完成，但有部分失败"
    echo "   请查看生成的测试报告了解详情"
fi

echo "=========================================="
echo "测试报告已生成到 api_test_report_*.json 文件" 