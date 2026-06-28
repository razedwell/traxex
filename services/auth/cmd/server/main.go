package main

import (
	"context"
	"log"

	"github.com/razedwell/traxex/shared/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Config loaded: %+v", cfg)

	select {
	case <-ctx.Done():
		log.Println("Context done")
		return
	default:
		log.Println("Context still running")
	}
}
