package main

import (
	"crud-go/internal/config"
	"crud-go/internal/repository/psql"
	"crud-go/internal/service"
	"crud-go/internal/transport/rest"
	"crud-go/pkg/database"
	"crud-go/pkg/hash"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func checkCurRelations(db *sql.DB) {
	// Query to retrieve all tables in the current database
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	var tables []string

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			logrus.Fatal(err)
		}
		tables = append(tables, table)
	}

	logrus.WithFields(logrus.Fields{
		"tables": tables,
	}).Info("Tables in the current database")

	if err = rows.Err(); err != nil {
		logrus.Fatal(err)
	}
}

func checkCurDB(db *sql.DB) {
	// Query to retrieve the name of the current database
	var dbName string
	err := db.QueryRow("SELECT current_database()").Scan(&dbName)
	if err != nil {
		logrus.Fatal(err)
	}

	// Print the name of the current database
	logrus.WithFields(logrus.Fields{
		"tables": dbName,
	}).Info("Current database")
}

// @title Phone API
// @description This is a RESTful API for managing phone records.
// @version 1.0
// @host localhost:8080
// @BasePath /api/phones

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {

	dbConfig, err := config.New()
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     dbConfig.DB.Host,
		Port:     dbConfig.DB.Port,
		Username: dbConfig.DB.Username,
		DBName:   dbConfig.DB.Name,
		SSLMode:  dbConfig.DB.SSLMode,
		Password: dbConfig.DB.Password,
	})

	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	checkCurRelations(db)
	checkCurDB(db)

	var b []byte

	phonesRepository := psql.NewPhone(db)
	usersRepository := psql.NewUser(db)
	phonesService := service.NewPhones(phonesRepository)
	usersService := service.NewUser(usersRepository, hash.NewSHA1Hasher("salt"), b, 2*time.Minute)
	controller := rest.NewController(phonesService, usersService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: controller.InitRouter(),
	}

	logrus.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}
}
