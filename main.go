package main

import (
	"context"
	"github.com/KodeKunstenaars/Share2Teach/db"
	"log"
)

func main() {
	// Connect to the database
	db.Connect()

	// Ensure disconnection from the database when the main function ends
	defer func() {
		if err := db.Client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
}
