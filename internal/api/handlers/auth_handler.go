package handlers

import (
	"blogSystem/internal/domain"
	"blogSystem/internal/service"
	"blogSystem/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	logger.Info("ShouldBind START...", zap.Any("c.PostForm:", c.HandlerName()))

	var req struct {
		Username string `json:"username" form:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" form:"password" binding:"required,min=3"`
		Email    string `json:"email" form:"email" binding:"required,email"`
	}

	if err := c.ShouldBind(&req); err != nil {
		logger.Info("ShouldBind START...", zap.Any("c.PostForm:", c.HandlerName()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := h.authService.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user_id": user.ID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful",
	})
}
