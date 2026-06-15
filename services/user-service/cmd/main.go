package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SovetkanB/payflow/user-service/internal/config"
	"github.com/SovetkanB/payflow/user-service/internal/db"
	"github.com/SovetkanB/payflow/user-service/internal/handler"
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

	db, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	log.Println("Connected to DB")

	h := handler.NewHandler()

	router := handler.NewRouter(h)

	log.Println("User Service starting on port", cfg.HTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
