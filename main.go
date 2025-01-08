package main

import (
	"BakeryService/db"
	"BakeryService/handlers"
	"BakeryService/logger"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var limiter = rate.NewLimiter(2, 5)

func main() {
	log := logger.InitLogger()

	logger.LogEntry(map[string]interface{}{
		"action":  "start",
		"status":  "success",
		"message": "Application started successfully",
		"time":    time.Now().Format(time.RFC3339),
	})

	database := db.ConnectDB()
	defer func() {
		sqlDB, _ := database.DB()
		sqlDB.Close()
		logger.LogEntry(map[string]interface{}{
			"action": "close_db",
			"status": "success",
			"time":   time.Now().Format(time.RFC3339),
		})
	}()

	r := mux.NewRouter()
	r.Use(rateLimitMiddleware)
	r.Use(errorHandlingMiddleware)

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/api/product/{id}", handlers.GetProductByID(database)).Methods("GET")
	r.HandleFunc("/api/product/create", handlers.CreateProduct(database)).Methods("POST")
	r.HandleFunc("/api/product/update", handlers.UpdateProduct(database)).Methods("PUT")
	r.HandleFunc("/api/product/delete/{id}", handlers.DeleteProduct(database)).Methods("DELETE")
	r.HandleFunc("/api/products", handlers.GetFilteredProducts(database)).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		logger.LogEntry(map[string]interface{}{
			"action": "shutdown",
			"status": "success",
			"time":   time.Now().Format(time.RFC3339),
		})
		if err := logger.FlushLogs("logs.json"); err != nil {
			log.WithError(err).Error("Failed to flush logs")
		}
		os.Exit(0)
	}()

	logger.LogEntry(map[string]interface{}{
		"address": "http://localhost:8080",
		"status":  "running",
		"time":    time.Now().Format(time.RFC3339),
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}

func errorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.LogEntry(map[string]interface{}{
					"error":    rec,
					"url":      r.URL.Path,
					"method":   r.Method,
					"clientIP": r.RemoteAddr,
					"headers":  r.Header,
					"time":     time.Now().Format(time.RFC3339),
				})

				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			w.Header().Set("X-RateLimit-Limit", "5")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Second).Format(time.RFC1123))

			logger.LogEntry(map[string]interface{}{
				"endpoint":    r.URL.Path,
				"clientIP":    r.RemoteAddr,
				"headers":     r.Header,
				"retry_after": time.Now().Add(time.Second).Format(time.RFC1123),
				"status":      "rate_limit_exceeded",
				"time":        time.Now().Format(time.RFC3339),
			})

			http.Error(w, `{"error":"Rate limit exceeded","message":"You have exceeded the allowed request rate. Please try again later."}`, http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Limit", "5")
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprint(limiter.Burst()-1))
		w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Second).Format(time.RFC1123))

		next.ServeHTTP(w, r)
	})
}
