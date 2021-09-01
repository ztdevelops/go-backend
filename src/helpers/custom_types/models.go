package custom_types

import (
	"gorm.io/gorm"
)

type User struct {
	ID       int
	Username string
	Password string
}

type Post struct {
	gorm.Model
	UserID   int
	User     User
	Message  string
	Likes    int
	Comments []Comment
}

type Comment struct {
	gorm.Model
	PostID  int
	Message string
	Likes   int
}