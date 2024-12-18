package handlers

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

func JSONHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		if r.Method == http.MethodPost {
			var data map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil || data["message"] == nil {
				json.NewEncoder(w).Encode(Response{
					Status:  "fail",
					Message: "Invalid JSON data",
				})
				return
			}

			fmt.Println("Message received:", data["message"])
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Data received successfully",
			})
		} else {
			json.NewEncoder(w).Encode(Response{
				Status:  "fail",
				Message: "Invalid request method",
			})
		}
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/menu.html")
}
