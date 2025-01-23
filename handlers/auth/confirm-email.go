package auth

import (
	"BakeryService/models"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

func ConfirmEmailHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Token is required", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("confirmation_token = ?", token).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Invalid or expired token", http.StatusBadRequest)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if user.EmailVerified {
			http.Error(w, "Email already verified", http.StatusBadRequest)
			return
		}

		err := db.Model(&user).Updates(map[string]interface{}{
			"email_verified":     true,
			"confirmation_token": nil,
		}).Error
		if err != nil {
			http.Error(w, "Failed to confirm email", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/static/email-confirmed.html", http.StatusSeeOther)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email confirmed successfully"})
	}
}
