package handlers

import (
    "distributed-file-storage/db"
    "distributed-file-storage/models"
    "github.com/gofiber/fiber/v2"
    "sync"
)

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