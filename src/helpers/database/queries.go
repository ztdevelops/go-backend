package database

import (
	"log"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
)

func (db *Database) SignUp(user custom_types.User) (err error) {
	if err = db.Model(&custom_types.User{}).Create(&user).Error; err != nil {
		return
	}
	log.Println("User successfully created:", user)
	return
}

func (db *Database) GetUser(username string) (user custom_types.User, err error) {
	if err = db.Model(&custom_types.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		return
	}
	return
}