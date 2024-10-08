package main

import (
	"distributed-file-storage/db"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connect to the database
	if err := db.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Fiber
	app := fiber.New()

	// Add routes here (we'll define these in later stages)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
