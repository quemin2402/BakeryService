package handlers

import (
	"BakeryService/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	"gorm.io/gorm"
)

func GetAllProducts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var products []models.Product
		if err := db.Order("id ASC").Find(&products).Error; err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func GetProductByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Error(w, "Product ID is required", http.StatusBadRequest)
			return
		}

		var product models.Product

		if err := db.First(&product, id).Error; err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

func CreateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		if product.Name == "" || product.Price <= 0 {
			http.Error(w, "Product name and price are required", http.StatusBadRequest)
			return
		}

		if err := db.Create(&product).Error; err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

func UpdateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Product ID is required", http.StatusBadRequest)
			return
		}

		var existingProduct models.Product
		if err := db.First(&existingProduct, id).Error; err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		var updatedProduct models.Product
		if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		if updatedProduct.Name != "" {
			existingProduct.Name = updatedProduct.Name
		}
		if updatedProduct.Price > 0 {
			existingProduct.Price = updatedProduct.Price
		}
		if updatedProduct.Description != "" {
			existingProduct.Description = updatedProduct.Description
		}

		if err := db.Save(&existingProduct).Error; err != nil {
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingProduct)
	}
}

func DeleteProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if err := db.Delete(&models.Product{}, id).Error; err != nil {
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
	}
}
