package handlers

import (
	"blogSystem/internal/domain"
	"blogSystem/internal/service"
	"blogSystem/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService1 *service.PostService) *PostHandler {
	return &PostHandler{postService: postService1}
}

// 创建文章
func (h *PostHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req struct {
		Title   string `json:"title" binding:"required,min=3,max=200"`
		Content string `json:"content" binding:"required,min=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &domain.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := h.postService.Create(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
		"user_id": post.UserID,
	})
}

// GetById 获取文章详情
func (h *PostHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	logger.Info("GetById START", zap.Error(err))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// 构建响应数据
	response := gin.H{
		"id":         post.ID,
		"title":      post.Title,
		"content":    post.Content,
		"user_id":    post.UserID,
		"created_at": post.CreatedAt,
		"author": gin.H{
			"id":       post.User.ID,
			"username": post.User.Username,
		},
	}
	// 添加评论数据
	if len(post.Comments) > 0 {
		var comments []gin.H
		for _, comment := range post.Comments {
			comments = append(comments, gin.H{
				"id":         comment.ID,
				"content":    comment.Content,
				"created_at": comment.CreatedAt,
				"user": gin.H{
					"id":       comment.User.ID,
					"username": comment.User.Username,
				},
			})
		}
		response["comments"] = comments
	}

	c.JSON(http.StatusOK, response)
}

// Update 更新文章
func (h *PostHandler) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req struct {
		Title   string `json:"title" binding:"omitempty,min=3,max=200"`
		Content string `json:"content" binding:"omitempty,min=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	if err := h.postService.Update(uint(id), userID, updates); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "update failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
}

// Delete 删除文章
func (h *PostHandler) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.postService.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "delete failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}

// List 获取文章列表
func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	posts, err := h.postService.List(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get posts"})
		return
	}

	var response []gin.H
	for _, post := range posts {
		response = append(response, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"created_at": post.CreatedAt,
			"author": gin.H{
				"id":       post.User.ID,
				"username": post.User.Username,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"page":  page,
		"size":  size,
		"total": len(response),
	})
}
