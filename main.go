package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"travel/internal/handler"
	"travel/internal/service"
	"travel/internal/storage"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	dbName := "travel"

	//create db connection
	db, err := createDatabase(dbName)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//create db migration
	m, err := setUpMigrations(db, dbName)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	//create storage
	storage := storage.New(db, "mysql")

	//create service
	service := service.New(storage)

	//create handler
	handler := handler.New(service)

	srv := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("I am running ... localhost:8080")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func createDatabase(dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(db:3306)/?parseTime=true&multiStatements=true")

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return nil, err
	}

	db.Close()

	db, err = sql.Open("mysql", "root:root@tcp(db:3306)/"+dbName+"?parseTime=true&multiStatements=true")

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func setUpMigrations(db *sql.DB, dbName string) (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(db, &mysql.Config{DatabaseName: dbName})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		fmt.Println("errorr")
		return nil, err
	}

	return m, nil
}
