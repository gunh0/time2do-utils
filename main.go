package main

import (
	"database/sql"
	"log"
	"os"
	"utils/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	config := database.Config{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		DB:       os.Getenv("MYSQL_DATABASE"),
	}

	connectionString := database.GetConnectionString(config)
	log.Println("[+] Connection String Check:", connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	} else {
		log.Println("Database connection successful")
	}
}
