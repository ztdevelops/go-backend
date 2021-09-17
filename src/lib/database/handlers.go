package database

import (
	"log"

	"github.com/ztdevelops/go-project/src/lib/custom"
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
	if err = database.MigrateTables(); err != nil {
		log.Println("failed to migrate tables:", err)
	}
}

func (database *Database) MigrateTables() (err error) {
	tables := []interface{}{
		custom.User{},
	}
	log.Println("Migrating tables.")

	for _, table := range tables {
		if err = database.AutoMigrate(table); err != nil {
			return
		}
	}
	log.Println("Migrate OK.")
	return
}
