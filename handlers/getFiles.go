package handlers

import (
	"distributed-file-storage/db"
	"distributed-file-storage/models"

	"github.com/gofiber/fiber/v2"
)

func GetFilesData(c *fiber.Ctx) error {
	fileID := c.Query("fileID")
	if fileID == "" {
		return c.Status(400).SendString("fileID is required")
	}

	var parts []models.FilePart
	db.DB.Where("file_id = ?", fileID).Order("part").Find(&parts)

	return c.JSON(parts)
}
