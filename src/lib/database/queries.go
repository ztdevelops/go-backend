package database

import (
	"log"

	"github.com/ztdevelops/go-project/src/lib/custom"
)

func (db *Database) SignUp(user custom.User) (err error) {
	if err = db.Model(&custom.User{}).Create(&user).Error; err != nil {
		return
	}
	log.Println("User successfully created:", user)
	return
}

func (db *Database) GetUser(username string) (user custom.User, err error) {
	if err = db.Model(&custom.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		return
	}
	return
}
