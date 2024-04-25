package main

import (
	"crud-go/internal/repository/psql"
	"crud-go/internal/service"
	"crud-go/internal/transport/rest"
	"crud-go/pkg/database"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func checkCurRelations(db *sql.DB) {
	// Query to retrieve all tables in the current database
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and print the table names
	log.Println("Tables in the current database:")
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tableName)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func checkCurDB(db *sql.DB) {
	// Query to retrieve the name of the current database
	var dbName string
	err := db.QueryRow("SELECT current_database()").Scan(&dbName)
	if err != nil {
		log.Fatal(err)
	}

	// Print the name of the current database
	log.Println("Current database:", dbName)
}

// @title Phone API
// @description This is a RESTful API for managing phone records.
// @version 1.0
// @host localhost:8080
// @BasePath /api/phones

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

	checkCurRelations(db)
	checkCurDB(db)

	phonesRepository := psql.NewPhone(db)
	phonesService := service.NewPhones(phonesRepository)
	phonesController := rest.NewPhonesHandler(phonesService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: phonesController.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
