package admin

import (
	"BakeryService/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func DeleteUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]

		if err := db.Delete(&models.User{}, userID).Error; err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	}
}

func UpdateUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]

		var req struct {
			Email     string `json:"email"`
			RoleID    uint   `json:"role_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		user.Email = req.Email
		user.RoleID = req.RoleID
		user.FirstName = req.FirstName
		user.LastName = req.LastName

		if err := db.Save(&user).Error; err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email     string `json:"email"`
			RoleID    uint   `json:"role_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("defaultpassword"), bcrypt.DefaultCost)

		user := models.User{
			Email:         req.Email,
			Password:      string(passwordHash),
			RoleID:        req.RoleID,
			EmailVerified: true,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
		}

		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
