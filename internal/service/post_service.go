package service

import (
	"blogSystem/internal/domain"
	"errors"

	"gorm.io/gorm"
)

// 4.文章管理功能
//
//	实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
//	实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
//	实现文章的更新功能，只有文章的作者才能更新自己的文章。
//	实现文章的删除功能，只有文章的作者才能删除自己的文章。
type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) Create(post *domain.Post) error {
	return s.db.Create(post).Error
}

func (s *PostService) GetByID(id uint) (*domain.Post, error) {
	var post domain.Post
	err := s.db.Preload("User").Preload("Comments.User").First(&post, id).Error
	return &post, err
}

func (s *PostService) Update(postID, userID uint, updates map[string]interface{}) error {
	return s.db.Debug().Model(&domain.Post{}).
		Where("id = ? AND user_id = ?", postID, userID).
		Updates(updates).Error
}

func (s *PostService) Delete(postID, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", postID, userID).Delete(&domain.Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("post not found or not owned by user")
	}
	return nil
}

func (s *PostService) List(page, size int) ([]domain.Post, error) {
	var posts []domain.Post
	err := s.db.Preload("User").
		Offset((page - 1) * size).
		Limit(size).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}
