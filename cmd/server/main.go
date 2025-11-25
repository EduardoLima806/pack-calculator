package main

import (
	"log"
	"net/http"

	"pack-calculator/internal/api"
	"pack-calculator/internal/config"
	"pack-calculator/internal/repository/memory"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	
	// Initialize repository with pack sizes from config
	packSizeRepo := memory.NewMemoryPackSizeRepository(cfg.PackSizes)
	
	// Create handler with repository
	handler := api.NewHandler(packSizeRepo)

	// Setup router
	r := mux.NewRouter()

	// API routes (registered first to take precedence)
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/calculate", handler.Calculate).Methods("POST")
	apiRouter.HandleFunc("/calculate", handler.CalculateQuery).Methods("GET")
	apiRouter.HandleFunc("/pack-sizes", handler.GetPackSizes).Methods("GET")

	// Health check
	r.HandleFunc("/health", handler.Health).Methods("GET")

	// Serve static files (UI) - catch-all for all other routes
	// This must be registered last so API routes are matched first
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	log.Printf("Default pack sizes: %v", cfg.PackSizes)
	
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

