package auth

import (
	"BakeryService/mailer"
	"github.com/google/uuid"
	"log"
	"net"
	"regexp"
	"strings"
)

import (
	"BakeryService/config"
	"BakeryService/models"
	"encoding/json"
	_ "fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func isEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return false
	}

	_, err := net.LookupMX(email[strings.Index(email, "@")+1:])
	return err == nil
}

func RegisterHandler(db *gorm.DB, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if !isEmailValid(req.Email) {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}

		confirmationToken := uuid.New().String()

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		newUser := models.User{
			Email:             req.Email,
			Password:          string(hashedPassword),
			RoleID:            2,
			EmailVerified:     false,
			ConfirmationToken: &confirmationToken,
		}

		if err := db.Create(&newUser).Error; err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		go func() {
			err := mailer.SendEmailWithConfirmation(config, req.Email, confirmationToken)
			if err != nil {
				log.Printf("Failed to send confirmation email: %v", err)
			}
		}()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered. Please check your email to verify your account."})
	}
}
