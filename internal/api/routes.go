package api

import (
	"blogSystem/internal/api/handlers"
	"blogSystem/internal/service"
	"blogSystem/pkg/auth"
	"blogSystem/pkg/database"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 获取数据库实例
	db := database.GetDB()

	// 初始化服务
	authService := service.NewAuthService(db)
	postService := service.NewPostService(db)
	commentService := service.NewCommentService(db)

	// 初始化服务器
	authHandler := handlers.NewAuthHandler(authService)

	// 公共路由
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// 需要认证的路由
	authGroup := r.Group("/")
	authGroup.Use(auth.JWTMiddleware())
	{
		// 文章路由
		postHandler := handlers.NewPostHandler(postService)
		authGroup.POST("/createPost", postHandler.Create)
		authGroup.GET("/getPostById/:id", postHandler.GetById)
		authGroup.POST("/UpdateById/:id", postHandler.Update)
		authGroup.GET("/DeleteById/:id", postHandler.Delete)
		authGroup.GET("/listPosts", postHandler.List)

		// 评论路由
		commentHandler := handlers.NewCommentHandler(commentService)
		authGroup.POST("/creatComment/:id", commentHandler.Create)
		authGroup.GET("/getCommentById/:id", commentHandler.GetByPostID)
		authGroup.GET("/deleteCommentById/:id", commentHandler.Delete)
	}

	return r
}
