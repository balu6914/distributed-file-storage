package main

import (
	"distributed-file-storage/db"
	"distributed-file-storage/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := db.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	app := fiber.New()

	app.Post("/upload", handlers.UploadFile)
	app.Get("/files", handlers.GetFilesData)
	app.Get("/download", handlers.DownloadFile)

	log.Fatal(app.Listen(":3000"))
}
