// docker build -t two-db .
// docker-compose up -d
package main

import (
	"CommentsService/db"
)

func main() {

	// Инициализация базы данных
	DB := db.InitDB()
	db.ExecuteSchemaSQL(DB)
	defer db.CloseDB()

}
