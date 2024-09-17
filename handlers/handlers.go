package handlers

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "log"
    "github.com/labstack/echo/v4"
    "github.com/arifrhm/go-upload-chunk/services"
    "github.com/arifrhm/go-upload-chunk/config"
    "github.com/arifrhm/go-upload-chunk/dto"
    "github.com/arifrhm/go-upload-chunk/middleware"
)

var (
    tempChunkSuffix = ".part"
)

func RegisterRoutes(e *echo.Echo) {
    e.GET("/upload", handleUploadPage)
    e.GET("/resume-upload", handleResumeUpload)
    e.POST("/upload-chunk", middleware.FileUploadMiddleware(handleUploadChunk))
}

// @title           File Upload API
// @version         1.0
// @description     API for handling file uploads in chunks.
// @host      localhost:8001
// @BasePath  /

// @Summary      Upload Page
// @Description  Serve the upload HTML page
// @Tags         upload
// @Produce      html
// @Success      200  {string}  string
// @Router       /upload [get]
func handleUploadPage(c echo.Context) error {
    htmlContent, err := os.ReadFile("upload.html")
    if err != nil {
        return c.String(
            http.StatusInternalServerError,
            "Failed to read HTML file!!!",
        )
    }
    return c.HTMLBlob(http.StatusOK, htmlContent)
}

// @Summary      Resume Upload
// @Description  Check uploaded chunks and return status
// @Tags         upload
// @Accept       json
// @Produce      json
// @Param        file_name  query   string  true  "File name"
// @Success      200  {object}  dto.ResumeUploadResponse
// @Router       /resume-upload [get]
func handleResumeUpload(c echo.Context) error {
    fileName := c.QueryParam("file_name")
    if fileName == "" {
        return c.JSON(
            http.StatusBadRequest,
            map[string]string{"message": "File name is required!!!"},
        )
    }

    chunks, err := filepath.Glob(
        filepath.Join(
            config.UploadPath,
            fmt.Sprintf("%s%s*", fileName, tempChunkSuffix),
        ),
    )
    if err != nil {
        return err
    }

    chunkIndex := len(chunks)
    totalChunks := chunkIndex // Assuming chunks are sequentially named

    return c.JSON(http.StatusOK, dto.ResumeUploadResponse{
        ChunkIndex:  chunkIndex,
        TotalChunks: totalChunks,
        FileName:    fileName,
    })
}

// @Summary      Upload Chunk
// @Description  Handle file chunk upload
// @Tags         upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        file      formData  file  true  "File chunk"
// @Param        chunk_index  formData  integer  true  "Chunk index"
// @Param        total_chunks  formData  integer  true  "Total chunks"
// @Param        file_name    formData  string  true  "File name"
// @Success      200  {object}  dto.UploadChunkResponse
// @Router       /upload-chunk [post]
func handleUploadChunk(c echo.Context) error {
    requestID := c.Get("request_id").(string)

    // Retrieve metadata from middleware context
    chunkIndex, ok := c.Get("chunkIndex").(int)
    if !ok {
        log.Printf("[ERROR] Request ID: %s - Invalid or missing chunkIndex", requestID)
        return c.JSON(http.StatusInternalServerError, "Internal server error")
    }

    totalChunks, ok := c.Get("totalChunks").(int)
    if !ok {
        log.Printf("[ERROR] Request ID: %s - Invalid or missing totalChunks", requestID)
        return c.JSON(http.StatusInternalServerError, "Internal server error")
    }

    file, err := c.FormFile("file")
    if err != nil {
        log.Printf("[ERROR] Request ID: %s - Failed to retrieve file: %v", requestID, err)
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "File is required"})
    }

    fileName := c.FormValue("file_name")
    if fileName == "" {
        log.Printf("[ERROR] Request ID: %s - File name is missing", requestID)
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "File name is required"})
    }

    // Log important details
    log.Printf("[INFO] Request ID: %s - Uploading chunk %d of %d for file %s", requestID, chunkIndex, totalChunks, fileName)

    // Open the file and handle the upload
    src, err := file.Open()
    if err != nil {
        log.Printf("[ERROR] Request ID: %s - Failed to open file: %v", requestID, err)
        return err
    }
    defer src.Close()

    if err := services.HandleFileUpload(src, file, fileName, chunkIndex, totalChunks); err != nil {
        log.Printf("[ERROR] Request ID: %s - Error in file upload: %v", requestID, err)
        return err
    }

    response := dto.UploadChunkResponse{
        RequestID: requestID,
    }

    // Check if all chunks are uploaded and assembled
    if chunkIndex == totalChunks-1 {
        response.Message = "All chunks uploaded and file assembled successfully!"
        log.Printf("[INFO] Request ID: %s - File %s successfully assembled", requestID, fileName)
    } else {
        response.Message = "Chunk uploaded successfully"
        log.Printf("[INFO] Request ID: %s - Chunk %d of %d uploaded successfully", requestID, chunkIndex, totalChunks)
    }

    return c.JSON(http.StatusOK, response)
}
