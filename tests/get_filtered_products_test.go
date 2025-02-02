package tests

import (
	"BakeryService/handlers"
	_ "BakeryService/handlers/admin"
	"BakeryService/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func InitTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	db.AutoMigrate(&models.Product{})

	return db
}

func TestGetFilteredProducts(t *testing.T) {
	db := InitTestDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	db.Create(&models.Product{ID: 1, Name: "French Baguette", Price: 3.00, Category: "Bread"})
	db.Create(&models.Product{ID: 2, Name: "Chocolate Croissant", Price: 3.50, Category: "Pastry"})

	req := httptest.NewRequest("GET", "/products?category=Bread", nil)
	w := httptest.NewRecorder()

	handler := handlers.GetFilteredProducts(db)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `[{"id":1,"name":"French Baguette","price":3.00,"description":"","image":"","category":"Bread"}]`
	assert.JSONEq(t, expected, w.Body.String())
}
