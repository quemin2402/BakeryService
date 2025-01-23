package db

import (
	"BakeryService/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	connStr := "host=localhost user=postgres dbname=bakery port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = db.AutoMigrate(&models.Email{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Connected to the database successfully!")

	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Order{}, &models.OrderProduct{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
