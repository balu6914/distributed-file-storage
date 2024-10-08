package handlers

import (
	"bytes"
	"distributed-file-storage/db"
	"distributed-file-storage/models"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func DownloadFile(c *fiber.Ctx) error {
	fileID := c.Query("fileID")
	if fileID == "" {
		return c.Status(400).SendString("fileID is required")
	}

	var parts []models.FilePart
	db.DB.Where("file_id = ?", fileID).Order("part").Find(&parts)

	var wg sync.WaitGroup
	var buffer bytes.Buffer

	for _, part := range parts {
		wg.Add(1)
		go func(data []byte) {
			defer wg.Done()
			buffer.Write(data)
		}(part.Data)
	}

	wg.Wait()

	return c.Send(buffer.Bytes())
}
