package auth

import (
	"BakeryService/models"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
)

type ProfileResponse struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthYear *int   `json:"birth_year"`
}

func ProfileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")[7:]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user models.User
		if err := db.Preload("Role").First(&user, claims.UserID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		response := ProfileResponse{
			Email:     user.Email,
			Role:      user.Role.Name,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			BirthYear: user.BirthYear,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
