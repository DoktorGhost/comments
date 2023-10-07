// docker build -t service-pg .
// docker run -d --name service-pg-container -p 5432:5432 service-pg
package main

import (
	"encoding/json"
	"log"
	"os"

	"GoNews/pcg/api"
	"GoNews/pcg/database"
	"GoNews/pcg/parse"
)

func main() {
	// Чтение конфигурационного файла
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to open config file:", err)
	}
	defer configFile.Close()

	var config parse.Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal("Failed to decode config file:", err)
	}

	// Инициализация базы данных
	db := database.InitDB()
	database.ExecuteSchemaSQL(db)
	defer db.Close()

	// Создание API сервера
	apiPort := "8080"
	go func() {
		err := api.StartAPI(apiPort, db)
		if err != nil {
			log.Fatal("Error starting API server:", err)
		}
	}()
}
