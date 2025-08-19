package main

import (
	"context"
	"log"

	"github.com/en7ka/chat-server/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err = a.Run(); err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
