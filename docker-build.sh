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

# 导出镜像到 tar 文件
docker save -o ${IMAGE_NAME}.tar ${IMAGE_NAME}:${TAG}

echo ""
echo "📦 正在上传 Docker 镜像到服务器..."

# 服务器配置
SERVER_HOST="154.219.110.51"
SERVER_PORT="8022"
SSH_KEY="~/.ssh/mmcc.ssh"
REMOTE_USER="mengchen"  # 可根据实际情况修改用户名

# 上传镜像文件到服务器
echo "📦 正在上传 Docker 镜像到服务器..."


# 上传镜像文件到服务器
sftp -P ${SERVER_PORT} -i ${SSH_KEY} ${REMOTE_USER}@${SERVER_HOST} << EOF
put ${IMAGE_NAME}.tar /tmp/${IMAGE_NAME}.tar
quit
EOF

if [ $? -eq 0 ]; then
    echo "✅ 镜像上传成功!"
else
    echo "❌ 镜像上传失败!"
    exit 1
fi

echo ""
echo "🚀 正在部署到服务器..."

# SSH 连接到服务器执行部署命令
ssh -p ${SERVER_PORT} -i ${SSH_KEY} ${REMOTE_USER}@${SERVER_HOST} << 'ENDSSH'
set -e

# 设置 sudo 密码
SUDO_PASSWORD="c5QDQvbh3uMdjpX3"

# 加载新镜像
echo "📦 正在加载 Docker 镜像..."
echo $SUDO_PASSWORD | sudo -S docker load -i /tmp/tangsong-esports.tar

# 停止并删除旧容器（如果存在）
echo "🧹 清理旧容器..."
if echo $SUDO_PASSWORD | sudo -S docker ps -q -f name=tangsong-esports-app; then
    echo "停止旧容器..."
    echo $SUDO_PASSWORD | sudo -S docker stop tangsong-esports-app
fi

if echo $SUDO_PASSWORD | sudo -S docker ps -aq -f name=tangsong-esports-app; then
    echo "删除旧容器..."
    echo $SUDO_PASSWORD | sudo -S docker rm tangsong-esports-app
fi

# 清理未使用的镜像
echo "🗑️ 清理未使用的镜像..."
echo $SUDO_PASSWORD | sudo -S docker image prune -f

# 启动新容器
echo "🚀 启动新容器..."
echo $SUDO_PASSWORD | sudo -S docker run -d --name tangsong-esports-app \
    -p 8000:8000 \
    -e DB_HOST=127.0.0.1 \
    -e DB_PORT=3306 \
    -e DB_USERNAME=tsshop \
    -e DB_PASSWORD=KEdG76bSGRhdWyhz \
    -e DB_NAME=tsshop \
    --restart unless-stopped \
    tangsong-esports:latest

# 检查容器状态
echo "📊 检查容器状态..."
echo $SUDO_PASSWORD | sudo -S docker ps -f name=tangsong-esports-app

# 等待容器完全启动
echo "⏳ 等待 15 秒让容器完全启动..."
sleep 15

# 检查容器日志
echo "📋 检查容器启动日志..."
echo $SUDO_PASSWORD | sudo -S docker logs tangsong-esports-app

# 清理临时文件
rm -f /tmp/tangsong-esports.tar

echo "✅ 部署完成!"
echo "🌐 服务访问地址: http://154.219.110.51:8000"
ENDSSH

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 部署成功完成!"
    echo "🌐 远程服务访问地址:"
    echo "   API: http://${SERVER_HOST}:8000/api/v1"
    echo "   Swagger: http://${SERVER_HOST}:8000/swagger/index.html"
else
    echo "❌ 部署失败!"
    exit 1
fi

# 清理本地临时文件
echo ""
echo "🧹 清理本地临时文件..."
rm -f ${IMAGE_NAME}.tar

echo ""
echo "✨ 所有操作完成!"


