package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:50;uniqueIndex;not null"`
	Password string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;uniqueIndex;not null"`
	Posts    []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Post struct {
	gorm.Model
	Title    string    `gorm:"size:200;not null"`
	Content  string    `gorm:"type:text;not null"`
	UserID   uint      `gorm:"index;not null"`
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"index;not null"`
	PostID  uint   `gorm:"index;not null"`
	User    User   `gorm:"foreignKey:UserID"`
	Post    Post   `gorm:"foreignKey:PostID"`
}
