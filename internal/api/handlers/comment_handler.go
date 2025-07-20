package handlers

import (
	"blogSystem/internal/domain"
	"blogSystem/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{service: s}
}

func (h *CommentHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Content string `json:"content" binding:"required,min=3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &domain.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}

	if err := h.service.Create(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      comment.ID,
		"content": comment.Content,
		"user_id": comment.UserID,
		"post_id": comment.PostID,
	})
}

func (h *CommentHandler) GetByPostID(c *gin.Context) {
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	comments, err := h.service.GetByPostID(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	commentID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.service.Delete(uint(commentID), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
