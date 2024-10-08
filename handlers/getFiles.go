package handlers

import (
	"distributed-file-storage/db"
	"distributed-file-storage/models"

	"github.com/gofiber/fiber/v2"
)

// GetFilesData godoc
// @Summary Retrieve file parts by file ID
// @Description Get all parts of a file by the provided file ID.
// @Tags file
// @Produce json
// @Param fileID query string true "Unique file ID"
// @Success 200 {array} models.FilePart
// @Failure 400 {string} string "fileID is required"
// @Router /files [get]
func GetFilesData(c *fiber.Ctx) error {
	fileID := c.Query("fileID")
	if fileID == "" {
		return c.Status(400).SendString("fileID is required")
	}

	var parts []models.FilePart
	db.DB.Where("file_id = ?", fileID).Order("part").Find(&parts)

	return c.JSON(parts)
}
