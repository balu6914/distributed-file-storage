package main

import (
	"log"

	"distributed-file-storage/db"
	"distributed-file-storage/handlers"

	_ "distributed-file-storage/docs" // Import the generated Swagger docs

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger" // Fiber middleware for Swagger
)

// @title Distributed File Storage API
// @version 1.0
// @description API for uploading, retrieving, and downloading files.
// @host localhost:3000
// @BasePath /
func main() {
	// Connect to the database
	if err := db.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Fiber app
	app := fiber.New()

	// Middleware for logging
	app.Use(logger.New())

	// Root route - Redirect to Swagger docs
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/")
	})

	// Swagger route
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	// Register API routes
	app.Post("/upload", handlers.UploadFile)    // Upload file API
	app.Get("/files", handlers.GetFilesData)    // Get uploaded file parts
	app.Get("/download", handlers.DownloadFile) // Download merged file

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
