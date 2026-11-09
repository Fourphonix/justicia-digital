package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
	})
}

func main() {
	mux := http.NewServeMux()
	
	// Ruta protegida por método GET nativo
	mux.HandleFunc("GET /health", healthHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Servidor Go corriendo en http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
