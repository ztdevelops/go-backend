package helpers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

)

func InitDatabaseConnection() {
	fmt.Println("Establishing connection with database.")
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
}