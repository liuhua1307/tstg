package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Mode     string
	Port     string
	Database DatabaseConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Charset  string
}

type JWTConfig struct {
	Secret string
	Expire int
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	viper.SetDefault("mode", "debug")
	viper.SetDefault("port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "root123456")
	viper.SetDefault("database.database", "tstg_shop")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("jwt.secret", "tangsong-esports-secret-key")
	viper.SetDefault("jwt.expire", 24)

	// 支持环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("配置文件读取失败，使用默认配置:", err)
	}

	AppConfig = &Config{
		Mode: viper.GetString("mode"),
		Port: viper.GetString("port"),
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			Database: viper.GetString("database.database"),
			Charset:  viper.GetString("database.charset"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("jwt.secret"),
			Expire: viper.GetInt("jwt.expire"),
		},
	}
}
