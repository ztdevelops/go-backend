package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (database *Database) InitDatabaseConnection() {
	log.Println("Establishing connection with database.")

	dsn := "root:password@tcp(127.0.0.1:3306)/local-production"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	log.Println("Database connection OK on port 127.0.0.1:3306")
	database.DB = db
	database.MigrateTables()
}

func (database *Database) MigrateTables() {
	log.Println("Migrating tables.")
}
