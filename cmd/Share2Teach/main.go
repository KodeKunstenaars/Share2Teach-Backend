package main

import (
	"context"
	"github.com/KodeKunstenaars/Share2Teach/internal/database"
	"github.com/KodeKunstenaars/Share2Teach/internal/router"
	"log"
)

func main() {
	// Connect to the database
	connect.Connect()

	// Set up the Gin router
	r := router.SetupRouter()
	if r != nil {
		err := r.Run(":8080")
		if err != nil {
			return
		}
	}

	// Ensure disconnection from the database when the main function ends
	defer func() {
		if err := connect.Client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
}
