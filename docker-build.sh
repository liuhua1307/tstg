#!/bin/bash

# 唐宋电竞陪玩报单平台 Docker 构建脚本

set -e

# 配置
IMAGE_NAME="tangsong-esports"
TAG="latest"
CONTAINER_NAME="tangsong-esports-app"

echo "🚀 开始构建唐宋电竞陪玩报单平台 Docker 镜像..."

# 构建镜像
docker build -t ${IMAGE_NAME}:${TAG} .

echo "✅ Docker 镜像构建完成!"
echo "📦 镜像名称: ${IMAGE_NAME}:${TAG}"

# 显示镜像信息
docker images | grep ${IMAGE_NAME}

echo ""
echo "🔧 使用方法:"
echo ""
echo "1. 快速启动 (使用默认配置):"
echo "   docker run -d --name ${CONTAINER_NAME} -p 8080:8080 ${IMAGE_NAME}:${TAG}"
echo ""
echo "2. 使用环境变量配置数据库:"
echo "   docker run -d --name ${CONTAINER_NAME} -p 8080:8080 \\"
echo "     -e DB_HOST=your_db_host \\"
echo "     -e DB_PORT=3306 \\"
echo "     -e DB_USERNAME=your_username \\"
echo "     -e DB_PASSWORD=your_password \\"
echo "     -e DB_NAME=your_database \\"
echo "     ${IMAGE_NAME}:${TAG}"
echo ""
echo "3. 使用自定义配置文件:"
echo "   docker run -d --name ${CONTAINER_NAME} -p 8080:8080 \\"
echo "     -v /path/to/your/config.yaml:/app/config.yaml \\"
echo "     ${IMAGE_NAME}:${TAG}"
echo ""
echo "4. 查看日志:"
echo "   docker logs -f ${CONTAINER_NAME}"
echo ""
echo "5. 停止容器:"
echo "   docker stop ${CONTAINER_NAME}"
echo ""
echo "6. 删除容器:"
echo "   docker rm ${CONTAINER_NAME}"
echo ""
echo "🌐 访问地址:"
echo "   API: http://localhost:8080/api/v1"
echo "   Swagger: http://localhost:8080/swagger/index.html" 