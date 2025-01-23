package admin

import (
	"BakeryService/models"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetUsersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		var users []models.User
		query := db.Preload("Role")
		if search != "" {
			query = query.Where("email ILIKE ?", "%"+search+"%")
		}

		if err := query.Find(&users).Error; err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetOrdersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		sort := r.URL.Query().Get("sort")
		order := r.URL.Query().Get("order")

		if order != "asc" && order != "desc" {
			order = "asc"
		}

		var orders []models.Order
		query := db.Preload("User").Preload("Products.Product")

		if search != "" {
			query = query.Joins("JOIN users ON users.id = orders.user_id").
				Where("users.email ILIKE ?", "%"+search+"%")
		}

		if sort == "order_date" {
			query = query.Order("order_date " + order)
		} else if sort == "total_price" {
			query = query.Order("total_price " + order)
		}

		if err := query.Find(&orders).Error; err != nil {
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}

func DeleteOrderHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Order ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}

		if err := db.Delete(&models.Order{}, id).Error; err != nil {
			http.Error(w, "Failed to delete order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Order deleted successfully"})
	}
}
