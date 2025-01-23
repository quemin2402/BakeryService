package db

import (
	"BakeryService/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	connStr := "host=dpg-cu98v31u0jms73fchiag-a.oregon-postgres.render.com user=bakery password=FOiJPHf6SQRqgUih0ipJN183iI1VA2mm dbname=bakery_688h port=5432 sslmode=require"
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
