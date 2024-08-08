package handlers

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/labstack/echo/v4"
    "github.com/arifrhm/go-upload-chunk/config"
    "github.com/arifrhm/go-upload-chunk/services"
)

var (
    tempChunkSuffix = ".part"
    uploadHTMLFile  = "./upload.html"
)

func UploadHandler(c echo.Context) error {
    htmlContent, err := os.ReadFile(uploadHTMLFile)
    if err != nil {
        return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read HTML file: %v", err))
    }
    return c.HTMLBlob(http.StatusOK, htmlContent)
}

func ResumeUploadHandler(c echo.Context) error {
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

func UploadChunkHandler(c echo.Context) error {
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
    chunkPath := filepath.Join(config.UploadPath, fmt.Sprintf("%s%s%d", fileName, tempChunkSuffix, chunkIndex))

    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(chunkPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    if _, err := io.Copy(dst, src); err != nil {
        return err
    }

    err = services.CheckAndAssembleFile(fileName, totalChunks)
    if err != nil {
        if err.Error() == "not all chunks are uploaded yet" {
            return c.JSON(http.StatusOK, map[string]string{"message": "Chunk uploaded successfully, awaiting more chunks"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error assembling file: %v", err)})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "File uploaded and assembled successfully!"})
}
