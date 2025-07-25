package main

import (
	"log"
	"net/http"
	"time"

	"github.com/sundaram2021/fruit-slot-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define routes
	r.GET("/play", handlers.Play)
	r.GET("/play/10", handlers.Play10)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Starting server on port 8080...")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}
