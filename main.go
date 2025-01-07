package main

import (
	"BakeryService/db"
	"BakeryService/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	database := db.ConnectDB()
	defer func(db *gorm.DB) {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}(database)

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/api/product/{id}", handlers.GetProductByID(database)).Methods("GET")
	r.HandleFunc("/api/product/create", handlers.CreateProduct(database)).Methods("POST")
	r.HandleFunc("/api/product/update", handlers.UpdateProduct(database)).Methods("PUT")
	r.HandleFunc("/api/product/delete/{id}", handlers.DeleteProduct(database)).Methods("DELETE")
	r.HandleFunc("/api/products", handlers.GetFilteredProducts(database)).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
