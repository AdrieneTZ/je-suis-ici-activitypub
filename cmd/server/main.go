package main

import (
	"fmt"
	"je-suis-ici-activitypub/internal/config"
	"je-suis-ici-activitypub/internal/db"
	"log"
)

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("fail to load config: %w", err)
	}

	// init database connection
	database, err := db.NewDatabase(db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("fail to connect database: %w", err)
	}
	defer database.Close()

	fmt.Println("FINALLY success connect to database!!!")
}
