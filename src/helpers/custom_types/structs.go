package custom_types

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

type Post struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey"`
	Message string `json:"message"`
	Username string `json:"username" gorm:"foreignKey:User.Username"`
	Likes int `json:"likes"`
}

type Comment struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey"`
	PostID []byte `json:"post_id" gorm:"foreignKey:Post.ID"`
	Message string `json:"message"`
	Username string `json:"username" gorm:"foreignKey:User.Username"`
	Likes int `json:"likes"`
}