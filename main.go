package main

import (
	"context"
	"github.com/gerhardotto/cmpg323_s2t/db"
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
