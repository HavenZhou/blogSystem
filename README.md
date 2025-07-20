# 个人博客系统

基于Gin+GORM实现的博客系统后端API

## 功能特性
用户认证：注册、登录（JWT认证）<br>
文章管理：创建、读取、分页获取文章列表、更新、删除文章<br>
评论功能：发表评论、获取文章评论列表<br>
错误处理：统一错误响应格式<br>
日志记录：请求日志和错误日志<br>

## 目录结构
```text
/blogSystem
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   ├── post_handler.go
│   │   │   └── comment_handler.go
│   │   └── routes.go
│   ├── domain/
│   │   └── models.go
│   └── service/
│       ├── auth_service.go
│       ├── post_service.go
│       └── comment_service.go
├── pkg/
│   ├── auth/
│   │   └── jwt.go
│   ├── database/
│   │   └── gorm.go
│   └── logger/
│       └── zap.go
├── .env
├── go.mod
├── go.sum
└── README.md
```

## 技术栈
语言：Go 1.16+<br>
Web框架：Gin<br>
ORM：GORM<br>
数据库：MySQL<br>
认证：JWT<br>
密码加密：bcrypt<br>
日志：ZAP<br>

## 快速开始

1. 配置环境变量
```env
DB_DSN="root:password@tcp(localhost:3306)/blog_test?charset=utf8mb4&parseTime=True"
JWT_SECRET="your-256-bit-secret"
SERVER_PORT="8080"
SERVER_ENV="development"
LOG_LEVEL="debug"
```
2. 启动服务
```env
go run cmd/main.go
```