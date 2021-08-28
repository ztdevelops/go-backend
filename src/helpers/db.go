package helpers

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabaseConnection() {
	fmt.Println("Establishing connection with database.")

	dsn := "root:password@tcp(127.0.0.1:3306)/test"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db2, _ := db.DB()
	defer db2.Close()
}