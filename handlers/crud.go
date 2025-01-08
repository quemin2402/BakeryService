package handlers

import (
	"BakeryService/logger"
	"BakeryService/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

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
			logger.LogEntry(map[string]interface{}{
				"endpoint": "GetProductByID",
				"status":   "failed",
				"error":    "Missing product ID",
				"time":     time.Now().Format(time.RFC3339),
			})
			http.Error(w, "Product ID is required", http.StatusBadRequest)
			return
		}

		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			logger.LogEntry(map[string]interface{}{
				"endpoint":   "GetProductByID",
				"product_id": id,
				"status":     "failed",
				"error":      err.Error(),
				"time":       time.Now().Format(time.RFC3339),
			})
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		logger.LogEntry(map[string]interface{}{
			"endpoint":   "GetProductByID",
			"product_id": id,
			"status":     "success",
			"time":       time.Now().Format(time.RFC3339),
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

func CreateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"endpoint": "CreateProduct",
				"error":    err,
			}).Error("Failed to parse request body")
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(product.Name) == "" {
			logrus.WithFields(logrus.Fields{
				"endpoint": "CreateProduct",
				"error":    "Product name is required",
			}).Warn("Product creation failed due to missing name")
			http.Error(w, "Product name is required", http.StatusBadRequest)
			return
		}

		if product.Price <= 0 {
			logrus.WithFields(logrus.Fields{
				"endpoint": "CreateProduct",
				"error":    "Invalid product price",
			}).Warn("Product creation failed due to invalid price")
			http.Error(w, "Price must be greater than zero", http.StatusBadRequest)
			return
		}

		if err := db.Create(&product).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"endpoint": "CreateProduct",
				"error":    err,
			}).Error("Failed to create product")
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		logrus.WithFields(logrus.Fields{
			"endpoint":   "CreateProduct",
			"product_id": product.ID,
		}).Info("Product created successfully")

		w.WriteHeader(http.StatusCreated)
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
			logger.LogEntry(map[string]interface{}{
				"endpoint":   "DeleteProduct",
				"product_id": id,
				"status":     "failed",
				"error":      err.Error(),
				"time":       time.Now().Format(time.RFC3339),
			})
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		logger.LogEntry(map[string]interface{}{
			"endpoint":   "DeleteProduct",
			"product_id": id,
			"status":     "success",
			"time":       time.Now().Format(time.RFC3339),
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
	}
}

func GetFilteredProducts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var products []models.Product

		nameFilter := r.URL.Query().Get("name")
		categoryFilter := r.URL.Query().Get("category")
		priceMin := r.URL.Query().Get("priceMin")
		priceMax := r.URL.Query().Get("priceMax")
		sortBy := r.URL.Query().Get("sortBy")
		sortOrder := r.URL.Query().Get("sortOrder")
		page := r.URL.Query().Get("page")

		limit := 8
		offset := 0
		if page != "" {
			pageNum, err := strconv.Atoi(page)
			if err == nil && pageNum > 1 {
				offset = (pageNum - 1) * limit
			}
		}

		query := db.Model(&models.Product{})

		if nameFilter != "" {
			query = query.Where("name ILIKE ?", "%"+nameFilter+"%")
		}
		if categoryFilter != "" {
			query = query.Where("category = ?", categoryFilter)
		}
		if priceMin != "" {
			if minPrice, err := strconv.ParseFloat(priceMin, 64); err == nil {
				query = query.Where("price >= ?", minPrice)
			}
		}
		if priceMax != "" {
			if maxPrice, err := strconv.ParseFloat(priceMax, 64); err == nil {
				query = query.Where("price <= ?", maxPrice)
			}
		}

		if sortBy != "" {
			order := "ASC"
			if sortOrder == "desc" {
				order = "DESC"
			}
			query = query.Order(sortBy + " " + order)
		}

		query = query.Limit(limit).Offset(offset)

		if err := query.Find(&products).Error; err != nil {
			logger.LogEntry(map[string]interface{}{
				"endpoint": "GetFilteredProducts",
				"filters": map[string]string{
					"name":      nameFilter,
					"category":  categoryFilter,
					"priceMin":  priceMin,
					"priceMax":  priceMax,
					"sortBy":    sortBy,
					"sortOrder": sortOrder,
				},
				"error":  err.Error(),
				"status": "failed",
				"time":   time.Now().Format(time.RFC3339),
			})
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}

		if len(products) == 0 {
			logger.LogEntry(map[string]interface{}{
				"endpoint": "GetFilteredProducts",
				"filters": map[string]string{
					"name":      nameFilter,
					"category":  categoryFilter,
					"priceMin":  priceMin,
					"priceMax":  priceMax,
					"sortBy":    sortBy,
					"sortOrder": sortOrder,
				},
				"product_count": 0,
				"status":        "no_results",
				"time":          time.Now().Format(time.RFC3339),
			})
			http.Error(w, "No products match the filter criteria", http.StatusNotFound)
			return
		}

		logger.LogEntry(map[string]interface{}{
			"endpoint": "GetFilteredProducts",
			"filters": map[string]string{
				"name":      nameFilter,
				"category":  categoryFilter,
				"priceMin":  priceMin,
				"priceMax":  priceMax,
				"sortBy":    sortBy,
				"sortOrder": sortOrder,
			},
			"product_count": len(products),
			"status":        "success",
			"time":          time.Now().Format(time.RFC3339),
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}
