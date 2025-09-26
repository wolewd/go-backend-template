package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-template/config"
)

func main() {
	// Initialize config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}
	defer cfg.DB.Close()

	log.Println("Config initialized successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down gracefully...")
}
