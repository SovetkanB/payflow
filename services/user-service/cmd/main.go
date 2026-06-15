package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SovetkanB/payflow/user-service/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found, using env variables: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("User Service starting on port", cfg.HTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
