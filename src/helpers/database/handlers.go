package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
)

func (database *Database) InitDatabaseConnection() {
	log.Println("Establishing connection with database.")

	dsn := "mysql://b6da1ad5ae2d23:019d058a@us-cdbr-east-04.cleardb.com"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	log.Println("Database connection OK/")
	database.DB = db
	database.MigrateTables()
}

func (database *Database) MigrateTables() {
	log.Println("Migrating tables.")
	database.AutoMigrate(&custom_types.User{})
}