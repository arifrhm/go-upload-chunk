package handlers

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "github.com/labstack/echo/v4"
    "github.com/arifrhm/go-upload-chunk/services"
    "github.com/arifrhm/go-upload-chunk/config"
)

var tempChunkSuffix = ".part"

func RegisterRoutes(e *echo.Echo) {
    e.GET("/upload", handleUploadPage)
    e.GET("/resume-upload", handleResumeUpload)
    e.POST("/upload-chunk", handleUploadChunk)
}

func handleUploadPage(c echo.Context) error {
    htmlContent, err := os.ReadFile("upload.html")
    if err != nil {
        return c.String(http.StatusInternalServerError, "Failed to read HTML file")
    }
    return c.HTMLBlob(http.StatusOK, htmlContent)
}

func handleResumeUpload(c echo.Context) error {
    fileName := c.QueryParam("file_name")
    if fileName == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "File name is required"})
    }

    chunks, err := filepath.Glob(filepath.Join(config.UploadPath, fmt.Sprintf("%s%s*", fileName, tempChunkSuffix)))
    if err != nil {
        return err
    }

    chunkIndex := len(chunks)
    totalChunks := chunkIndex // Assuming chunks are sequentially named

    return c.JSON(http.StatusOK, map[string]interface{}{
        "chunk_index":  chunkIndex,
        "total_chunks": totalChunks,
        "file_name":    fileName,
    })
}

func handleUploadChunk(c echo.Context) error {
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }

    chunkIndex, err := strconv.Atoi(c.FormValue("chunk_index"))
    if err != nil {
        return err
    }

    totalChunks, err := strconv.Atoi(c.FormValue("total_chunks"))
    if err != nil {
        return err
    }

    fileName := c.FormValue("file_name")
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    if err := services.HandleFileUpload(src, file, fileName, chunkIndex, totalChunks); err != nil {
        return err
    }

    // Check if all chunks have been uploaded and assembled
    chunks, err := filepath.Glob(filepath.Join(config.UploadPath, fmt.Sprintf("%s%s*", fileName, tempChunkSuffix)))
    if err != nil {
        return err
    }

    if len(chunks) == 0 {
        return c.JSON(http.StatusOK, map[string]string{"message": "All chunks uploaded and file assembled successfully!"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Chunk uploaded successfully"})
}
