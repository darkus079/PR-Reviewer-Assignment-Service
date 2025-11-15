package main

import (
	"log"

	"pr-reviewer-assignment-service/internal/config"
	"pr-reviewer-assignment-service/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Server starting on port", cfg.Server.Port)
}

