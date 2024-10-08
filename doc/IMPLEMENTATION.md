# Implementation Details

## Overview

This document provides a detailed overview of the implementation of the **Distributed File Storage** project. It includes explanations of the architecture, data model, key features, and multithreading approach used to handle file uploads, retrieval, and downloads. The project is built with Golang and is designed to be scalable, efficient, and easy to set up using Docker.

## Architecture

The project follows a modular architecture with the following components:

- **Database Layer (`db`)**: Handles the connection to the SQLite database using GORM, an ORM library for Golang.
- **API Layer (`handlers`)**: Contains route handlers for handling API requests for file uploads, retrieval, and downloads.
- **Models (`models`)**: Defines the data structure used to store file parts.
- **Swagger Documentation (`docs`)**: Auto-generated documentation for the API using `swag`.

## Data Model

### FilePart Model

The `FilePart` model is used to store each part of an uploaded file. It is defined in `models/file.go`:

```go
type FilePart struct {
    ID     uint   `gorm:"primaryKey"`
    FileID string `gorm:"index"`
    Part   int
    Data   []byte
}
```

- **ID**: Primary key for each file part (auto-incremented).
- **FileID**: A unique identifier for the entire file.
- **Part**: The part number of the file, used to maintain the order of parts.
- **Data**: The binary data of the file part.

## Key Features

### 1. File Upload

- **Handler**: `handlers/upload.go`
- **Function**: `UploadFile`
- **Description**: The upload functionality is designed to handle large files by splitting them into smaller parts (default size: 1MB). Each part is stored in the database as a separate record.
- **Multithreading**: 
  - Uses Golang's `sync.WaitGroup` to manage parallel storage of file parts.
  - Each file part is saved using a separate Goroutine, enabling concurrent database writes.

```go
wg.Add(1)
go func(part int, data []byte) {
    defer wg.Done()
    db.DB.Create(&models.FilePart{
        FileID: fileID,
        Part:   part,
        Data:   data,
    })
}(part, partData)
```

- **Database**: Stores each part in the SQLite database using GORM.

### 2. File Retrieval

- **Handler**: `handlers/getFiles.go`
- **Function**: `GetFilesData`
- **Description**: Retrieves all parts of a file based on the provided `fileID`. The file parts are sorted by their `Part` number to ensure correct ordering when the file is reconstructed.
- **Response**: Returns a JSON array of `FilePart` objects.

```go
db.DB.Where("file_id = ?", fileID).Order("part").Find(&parts)
```

### 3. File Download

- **Handler**: `handlers/download.go`
- **Function**: `DownloadFile`
- **Description**: Merges all parts of a file based on `fileID` and returns the complete file as a binary stream. This allows users to download the reconstructed file.
- **Multithreading**: 
  - Uses `sync.WaitGroup` and a `bytes.Buffer` to merge file parts in parallel.
  - Each part is written to a buffer using a Goroutine.

```go
var buffer bytes.Buffer
for _, part := range parts {
    wg.Add(1)
    go func(data []byte) {
        defer wg.Done()
        buffer.Write(data)
    }(part.Data)
}
```

### 4. API Documentation

- **Tool**: `swaggo/swag`
- **Setup**: Uses the `swag init` command to generate Swagger documentation from comments in the code.
- **Access**: The documentation is available at `http://localhost:3000/docs/` when the server is running.
- **Benefits**: Enables users to interact with the API directly from their browsers and provides detailed information about each endpoint, making testing easier.

## Multithreading Approach

### Why Multithreading?

- **Efficiency**: By using multithreading, the server can handle file uploads and downloads more efficiently, especially for large files.
- **Concurrency**: Each part of a file is processed concurrently, reducing the time needed to store or retrieve files.
- **Scalability**: This approach allows the server to handle multiple upload and download requests simultaneously, making it more scalable.

### How It Works

- **Upload**: Each file part is stored using a separate Goroutine, allowing parts to be saved in parallel.
- **Download**: During the download, each part is written to a buffer using a separate Goroutine, ensuring quick reconstruction of the file.
- **Sync Mechanism**: Uses `sync.WaitGroup` to ensure that all parts are saved or merged before sending a response back to the user.

## Dockerization

- **Dockerfile**: Defines a container for running the Golang server. It includes building the application and running it.
- **docker-compose.yml**: Simplifies the process of building and running the container using `docker-compose up --build`.
- **Benefits**: 
  - Makes the setup process easier, as all dependencies are packaged into a container.
  - Ensures that the application runs consistently across different environments.

### Docker Commands

- **Build and Run**:
   ```bash
   docker-compose up --build
   ```

- **Stop and Remove Containers**:
   ```bash
   docker-compose down
   ```

## Error Handling

- **Validation**: Ensures that required parameters like `fileID` are provided for the `/files` and `/download` endpoints.
- **Graceful Error Responses**: Returns appropriate status codes (`400`, `500`, etc.) with meaningful messages for client-side and server-side errors.
- **Golang Error Handling**: Uses `if err != nil` checks throughout the code to handle potential errors and return meaningful feedback.

## Future Improvements

- **Authentication**: Add user authentication to secure access to the API.
- **File Encryption**: Implement encryption for file data to enhance security.
- **Support for Additional Databases**: Expand support for other databases like PostgreSQL or MySQL.
- **Performance Optimization**: Fine-tune the multithreading approach for handling extremely large files.

## Conclusion

The **Distributed File Storage** project is designed to handle large file uploads and downloads efficiently using Golang’s concurrency features. With its modular structure, it is easy to maintain and extend. The use of Swagger makes it user-friendly for developers who want to understand and test the API. Dockerization ensures that the setup process is simple, making it easy to deploy the application in different environments.
