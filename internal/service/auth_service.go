package service

import (
	"blogSystem/internal/domain"
	"blogSystem/pkg/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(user *domain.User) error {
	// 检查用户名是否已存在
	var count int64
	s.db.Model(&domain.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("username already exists")
	}

	// 密码加密
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	return s.db.Create(user).Error
}

func (s *AuthService) Login(username, password string) (string, error) {
	var user domain.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateToken(user.ID)
}
