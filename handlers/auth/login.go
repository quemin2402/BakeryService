package auth

import (
	"BakeryService/config"
	"BakeryService/models"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var JwtKey = []byte("secret_key")

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func LoginHandler(db *gorm.DB, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		var userID uint
		var role string

		if req.Email == config.AdminEmail && req.Password == config.AdminPassword {
			role = "admin"
			userID = 0
		} else {
			var user models.User
			if err := db.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}

			if !user.EmailVerified {
				http.Error(w, "Please verify your email before logging in", http.StatusUnauthorized)
				return
			}

			role = user.Role.Name
			userID = user.ID
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			UserID: userID,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
			Path:    "/",
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Login successful",
			"token":   tokenString,
		})
	}
}
