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
var broadcast = make(chan models.MessageResponse)
var mutex sync.Mutex

func HandleConnections(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	chatIDStr := r.URL.Query().Get("chat_id")

	if chatIDStr == "" || chatIDStr == "null" {
		log.Println("Invalid chat ID received:", chatIDStr)
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil || chatID <= 0 {
		log.Println("Invalid chat ID format:", chatIDStr)
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var userID uint
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Println("JWT validation failed:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID = claims.UserID

	log.Printf("User %d connecting to chat %d", userID, chatID)

	var chat models.Chat
	if err := db.First(&chat, chatID).Error; err != nil {
		log.Println("Chat not found:", chatID)
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	clients[userID] = ws
	mutex.Unlock()

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading JSON:", err)
			mutex.Lock()
			delete(clients, userID)
			mutex.Unlock()
			break
		}

		msg.ChatID = uint(chatID)
		msg.Timestamp = time.Now()

		msg.UserID = userID
		if err := db.Create(&msg).Error; err != nil {
			log.Println("Error saving message:", err)
		}

		var role uint

		if claims.Role == "admin" {
			role = 1
		} else {
			role = 2
		}

		responseMsg := models.MessageResponse{
			ID:        msg.ID,
			Content:   msg.Content,
			UserID:    msg.UserID,
			ChatID:    msg.ChatID,
			Timestamp: msg.Timestamp,
			Role:      role,
		}

		broadcast <- responseMsg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		log.Printf("Broadcasting message: %+v", msg)

		mutex.Lock()
		for userID, client := range clients {

			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error sending message to user", userID, ":", err)
				client.Close()
				delete(clients, userID)
			}

		}
		mutex.Unlock()
	}
}

func StartChat(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) < 7 {
			log.Println("Invalid token received")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		token = token[7:]

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			log.Println("JWT validation failed:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims.UserID
		log.Printf("User %d is trying to start a chat", userID)

		// Simply create a new chat
		chat := models.Chat{
			UserID:   userID,
			IsActive: true,
		}

		if err := db.Create(&chat).Error; err != nil {
			log.Println("Failed to create chat:", err)
			http.Error(w, "Failed to create chat", http.StatusInternalServerError)
			return
		}

		log.Printf("New chat created: ChatID %d for User %d", chat.ID, userID)

		response := map[string]interface{}{"id": chat.ID}
		log.Printf("Sending new chat data: %+v", response)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chat)
	}
}

func CloseChat(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIDStr := r.URL.Query().Get("chat_id")
		if chatIDStr == "" {
			http.Error(w, "Missing chat ID", http.StatusBadRequest)
			return
		}

		chatID, err := strconv.Atoi(chatIDStr)
		if err != nil || chatID <= 0 {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		token := r.Header.Get("Authorization")
		if len(token) < 7 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(token[7:])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			http.Error(w, "Forbidden: Only admin can close chats", http.StatusForbidden)
			return
		}

		var chat models.Chat
		if err := db.First(&chat, chatID).Error; err != nil {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}

		if err := db.Model(&models.Chat{}).Where("id = ?", chatID).Update("is_active", false).Error; err != nil {
			http.Error(w, "Failed to close chat", http.StatusInternalServerError)
			return
		}

		mutex.Lock()
		if client, ok := clients[chat.UserID]; ok {
			client.Close()
			delete(clients, chat.UserID)
		}
		mutex.Unlock()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Chat closed successfully"})
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

func GetChatMessages(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := r.URL.Query().Get("chat_id")
		if chatID == "" {
			http.Error(w, "Missing chat ID", http.StatusBadRequest)
			return
		}

		token := r.Header.Get("Authorization")
		if len(token) < 7 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ValidateJWT(token[7:])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		cid, err := strconv.Atoi(chatID)
		if err != nil || cid <= 0 {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		chatIDUint := uint(cid)

		var chat models.Chat
		if err := db.First(&chat, chatIDUint).Error; err != nil {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}

		if claims.Role != "admin" && chat.UserID != claims.UserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var messages []struct {
			models.Message
			RoleID uint `json:"role_id"`
		}

		if err := db.Table("messages").
			Select("messages.*, users.role_id").
			Joins("LEFT JOIN users ON messages.user_id = users.id").
			Where("messages.chat_id = ?", chatIDUint).
			Order("messages.timestamp asc").
			Scan(&messages).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		var responseMessages []models.MessageResponse
		for _, msg := range messages {
			roleID := msg.RoleID
			if msg.UserID == 0 {
				roleID = 1
			}

			responseMsg := models.MessageResponse{
				ID:        msg.ID,
				Content:   msg.Content,
				UserID:    msg.UserID,
				ChatID:    msg.ChatID,
				Timestamp: msg.Timestamp,
				Role:      roleID,
			}
			responseMessages = append(responseMessages, responseMsg)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseMessages)
	}
}
