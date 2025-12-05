package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB menginisialisasi koneksi database PostgreSQL
func InitDB() {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "student"
	dbname := "mkp_ticketing"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Test koneksi
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("Successfully connected to database!")
}

// CloseDB menutup koneksi database
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
