package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDB() (*Database, error) {
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
		
	dsn := fmt.Sprintf("postgresql://%s:%s@%s-8649.8nk.gcp-asia-southeast1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full", user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		fmt.Println("Successfully connected to the database")
	}


	return &Database {
		DB:   db,
	}, nil
}