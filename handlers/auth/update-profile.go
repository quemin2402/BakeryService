package auth

import (
	"BakeryService/models"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type UpdateProfileRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	BirthYear *int   `json:"birth_year,omitempty"`
}

func UpdateProfileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req UpdateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			log.Printf("Error finding user with ID %d: %v", userID, err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if req.FirstName != "" {
			user.FirstName = req.FirstName
		}
		if req.LastName != "" {
			user.LastName = req.LastName
		}
		if req.BirthYear != nil {
			user.BirthYear = req.BirthYear
		}

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user with ID %d: %v", userID, err)
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
	}
}
