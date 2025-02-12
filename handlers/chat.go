package handlers

import (
	"BakeryService/handlers/auth"
	"BakeryService/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[uint]*websocket.Conn)
var broadcast = make(chan models.Message)
var mutex sync.Mutex

func HandleConnections(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	chatIDStr := r.URL.Query().Get("chat_id")

	if chatIDStr == "" {
		log.Println("Invalid chat ID: null")
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil || chatID <= 0 {
		log.Println("Invalid chat ID:", chatIDStr)
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var userID int
	if token == "" {
		log.Println("Super Admin detected, assigning userID = -1")
		userID = 0 // Фиктивный ID для админа без токена
	} else {
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			log.Println("JWT validation failed:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID = int(claims.UserID)
	}

	// Проверяем, существует ли чат
	var chat models.Chat
	if err := db.First(&chat, chatID).Error; err != nil {
		log.Println("Chat not found:", chatID)
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	log.Printf("User %d connected to WebSocket in chat %d", userID, chatID)

	mutex.Lock()
	clients[uint(userID)] = ws
	mutex.Unlock()

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading JSON:", err)
			mutex.Lock()
			delete(clients, uint(userID))
			mutex.Unlock()
			break
		}

		// Проверяем, можно ли сохранять сообщения
		if userID > 0 {
			msg.UserID = uint(userID)
			msg.ChatID = uint(chatID)
			msg.Timestamp = time.Now()

			log.Printf("Processed message: %+v", msg)

			if err := db.Create(&msg).Error; err != nil {
				log.Println("Error saving message:", err)
			}
		} else {
			log.Println("Admin is sending a message, skipping database insert")
		}

		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		log.Printf("Broadcasting message: %+v", msg)

		mutex.Lock()
		for userID, client := range clients {
			if msg.UserID != userID {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Println("Error sending message to user", userID, ":", err)
					client.Close()
					delete(clients, userID)
				}
			}
		}
		mutex.Unlock()
	}
}

func StartChat(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) < 7 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		token = token[7:]

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims.UserID
		if userID == 0 {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		log.Printf("User %d is trying to start a chat", userID)

		var existingChat models.Chat
		if err := db.Where("user_id = ? AND is_active = true", userID).First(&existingChat).Error; err == nil {
			log.Printf("Chat already exists for user %d: ChatID %d", userID, existingChat.ID)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(existingChat)
			return
		}

		chat := models.Chat{UserID: userID, IsActive: true}
		if err := db.Create(&chat).Error; err != nil {
			log.Println("Failed to create chat:", err)
			http.Error(w, "Failed to create chat", http.StatusInternalServerError)
			return
		}

		log.Printf("New chat created: ChatID %d for User %d", chat.ID, userID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chat)
	}
}

func CloseChat(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) < 7 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		token = token[7:]

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims.UserID

		db.Model(&models.Chat{}).Where("user_id = ?", userID).Update("is_active", false)

		mutex.Lock()
		if client, ok := clients[userID]; ok {
			client.Close()
			delete(clients, userID)
		}
		mutex.Unlock()

		w.WriteHeader(http.StatusOK)
	}
}

func GetActiveChats(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		var isAdmin bool
		var _ uint

		if len(token) < 7 {
			isAdmin = true
		} else {
			token = token[7:]

			claims, err := auth.ValidateJWT(token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_ = claims.UserID
			isAdmin = claims.Role == "admin"
		}

		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var chats []models.Chat
		if err := db.Where("is_active = ?", true).Find(&chats).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		log.Printf("Sending active chats: %+v", chats)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chats)
	}
}

func GetChatStatus(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) < 7 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		token = token[7:]

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims.UserID
		if userID == 0 {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		var chat models.Chat
		if err := db.Where("user_id = ? AND is_active = true", userID).First(&chat).Error; err == nil {
			json.NewEncoder(w).Encode(chat)
			return
		}

		http.Error(w, "No active chat", http.StatusNotFound)
	}
}
