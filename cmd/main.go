package main

import (
	"blogSystem/config"
	"blogSystem/internal/api"
	"blogSystem/pkg/database"
	"blogSystem/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 第一阶段：基础日志初始化（用于配置加载期间的日志）
	logger.Init("info")
	defer logger.Sync()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 第二阶段：用配置中的设置重新初始化日志
	logger.Init(cfg.Log.Level)
	defer logger.Sync()

	// 初始化数据库
	if err := database.Init(cfg.DB.DSN); err != nil {
		logger.Fatal("Database initialization failed",
			zap.String("dsn", cfg.DB.DSN),
			zap.Error(err),
		)
	}
	defer database.Close()

	// 初始化HTTP服务器
	router := api.NewRouter()
	logger.Info("Server is starting",
		zap.String("port", cfg.Server.Port),
		zap.String("log_level", cfg.Log.Level),
	)

	// 启动服务器
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Server startup failed",
			zap.String("port", cfg.Server.Port),
			zap.Error(err),
		)
	}

	// 服务关闭
	logger.Fatal("Server stopped", zap.Error(router.Run(":"+cfg.Server.Port)))
}
