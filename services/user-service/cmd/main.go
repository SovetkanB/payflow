package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SovetkanB/payflow/user-service/internal/config"
	"github.com/SovetkanB/payflow/user-service/internal/db"
	"github.com/SovetkanB/payflow/user-service/internal/handler"
	"github.com/SovetkanB/payflow/user-service/internal/repository"
	"github.com/SovetkanB/payflow/user-service/internal/service"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("No .env file found, using env variables: %v", err)
	// }

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Close()
	log.Println("Connected to DB")

	if err := db.RunMigrate(cfg.MigrationDSN(), cfg.MigrationsPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied successfully")

	r := repository.NewRepository(database)
	s := service.NewService(r, cfg)
	h := handler.NewHandler(s)

	router := handler.NewRouter(h, cfg)

	log.Println("User Service starting on port", cfg.HTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
