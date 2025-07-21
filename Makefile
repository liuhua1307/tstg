.PHONY: run build clean test deps swagger

# 默认目标
all: deps build

# 安装依赖
deps:
	go mod tidy
	go mod download

# 运行开发服务器
run:
	go run main.go

# 构建可执行文件
build:
	go build -o bin/tangsong-esports main.go

# 构建生产版本
build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/tangsong-esports main.go

# 清理构建文件
clean:
	rm -rf bin/

# 运行测试
test:
	go test -v ./...

# 生成 Swagger 文档
swagger:
	swag init

# 格式化代码
fmt:
	go fmt ./...

# 检查代码
lint:
	golangci-lint run

# 启动数据库（使用 Docker）
db-up:
	docker run --name tangsong-mysql \
		-e MYSQL_ROOT_PASSWORD=123456 \
		-e MYSQL_DATABASE=tangsong_esports \
		-p 3306:3306 \
		-d mysql:8.0

# 停止数据库
db-down:
	docker stop tangsong-mysql && docker rm tangsong-mysql

# 数据库迁移
migrate:
	go run main.go

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  deps        - Install dependencies"
	@echo "  run         - Run development server"
	@echo "  build       - Build executable"
	@echo "  build-prod  - Build production executable"
	@echo "  clean       - Clean build files"
	@echo "  test        - Run tests"
	@echo "  swagger     - Generate Swagger docs"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  db-up       - Start MySQL with Docker"
	@echo "  db-down     - Stop MySQL Docker container"
	@echo "  migrate     - Run database migration"
	@echo "  help        - Show this help" 