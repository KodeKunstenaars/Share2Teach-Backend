package main

import (
	"github.com/KodeKunstenaars/Share2Teach/internal/router"
)

func main() {

	// Set up the Gin router
	r := router.SetupRouter()
	if r != nil {
		err := r.Run(":8080")
		if err != nil {
			return
		}
	}

}
