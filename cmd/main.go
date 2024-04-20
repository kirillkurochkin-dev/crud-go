package main

import (
	"crud-go/pkg/database"
	"log"
	"net/http"
)

func main() {
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "crud-go",
		SSLMode:  "disable",
		Password: "postgres",
	})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//TODO: Инициализировать зависимости (репозиторий -> сервис -> хендлер)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
