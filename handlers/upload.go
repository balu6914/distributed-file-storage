package handlers

import (
	"distributed-file-storage/db"
	"distributed-file-storage/models"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// UploadFile godoc
// @Summary Upload a file and split it into parts
// @Description Upload a file, split it into parts, and store each part in the database.
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param fileID formData string true "Unique file ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "File not provided"
// @Failure 500 {string} string "Failed to open file"
// @Router /upload [post]
func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("File not provided")
	}

	fileID := c.FormValue("fileID") // Generate or provide a unique file ID.

	// Open the file for reading.
	f, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Failed to open file")
	}
	defer f.Close()

	// Read the file and divide it into parts.
	const partSize = 1024 * 1024 // 1MB parts
	var part int
	buffer := make([]byte, partSize)
	var wg sync.WaitGroup

	for {
		n, err := f.Read(buffer)
		if n == 0 || err != nil {
			break
		}

		partData := make([]byte, n)
		copy(partData, buffer[:n])

		wg.Add(1)
		go func(p int, data []byte) {
			defer wg.Done()
			db.DB.Create(&models.FilePart{
				FileID: fileID,
				Part:   p,
				Data:   data,
			})
		}(part, partData)

		part++
	}

	wg.Wait()
	return c.JSON(fiber.Map{"fileID": fileID})
}
