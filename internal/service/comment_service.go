package service

import (
	"blogSystem/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) Create(comment *domain.Comment) error {
	return s.db.Create(comment).Error
}

func (s *CommentService) GetByPostID(postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := s.db.Debug().Preload("User").Where("post_id = ?", postID).Find(&comments).Error
	return comments, err
}

func (s *CommentService) Delete(commentID, userID uint) error {
	result := s.db.Debug().Where("id = ? AND user_id = ?", commentID, userID).Delete(&domain.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("comment not found or not owned by user")
	}
	return nil
}
