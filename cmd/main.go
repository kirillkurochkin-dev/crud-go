package main

import (
	"crud-go/pkg/database"
	"database/sql"
	"log"
	"net/http"
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

	//TODO: Инициализировать зависимости (репозиторий -> сервис -> хендлер)
	//http.HandleFunc("/phones", handleFunc)
	//http.HandleFunc("/phonesGet", handleFuncGet)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
