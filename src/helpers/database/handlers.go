package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
)

func (database *Database) InitDatabaseConnection() (err error) {
	log.Println("Establishing connection with database.")

	dsn := "root:password@tcp(127.0.0.1:3306)/backend-main"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	log.Println("Database connection OK on port 127.0.0.1:3306")
	database.DB = db
	if err = database.MigrateTables(); err != nil {
		return
	}
	return
}

func (database *Database) MigrateTables() (err error) {
	log.Println("Migrating tables.")

	tablesToMigrate := []interface{}{
		custom_types.User{},
		custom_types.Post{},
		custom_types.Comment{},
	}

	for _, table := range tablesToMigrate {
		tableName := custom_types.GetStructName(table)
		if err = database.AutoMigrate(table); err != nil {
			log.Println(tableName, "failed to migrate.")
			return
		}
		log.Println(tableName, "OK.")
	}

	log.Println("Migration OK.")
	return
}
