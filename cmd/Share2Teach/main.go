package main

import (
	"context"
	"github.com/KodeKunstenaars/Share2Teach/internal/database"
	"log"
)

func main() {
	// Connect to the database
	connect.Connect()

	// Ensure disconnection from the database when the main function ends
	defer func() {
		if err := connect.Client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
}
