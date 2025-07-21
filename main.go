package main

import (
	"log"
	"tangsong-esports/config"
	"tangsong-esports/database"
	"tangsong-esports/router"
	"tangsong-esports/utils"

	"github.com/gin-gonic/gin"
)

// @title 唐宋电竞陪玩报单平台API
// @version 1.0
// @description 唐宋电竞陪玩报单平台后端API文档
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	database.InitDB()

	// 初始化日志
	utils.InitLogger()

	// 设置Gin模式
	if config.AppConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化路由
	r := router.InitRouter()

	// 启动服务器
	port := config.AppConfig.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
