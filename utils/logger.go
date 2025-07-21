package utils

import (
	"log"
	"os"
)

func InitLogger() {
	// 设置日志输出格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
	log.Println("日志系统初始化完成")
}
