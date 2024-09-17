package middleware

import (
    "bytes"
    "encoding/hex"
    "io"
    "log"
    "net/http"
    "strconv"
    "github.com/labstack/echo/v4"
)

const (
    maxChunkSize = 1 * 1024 * 1024 // 1 MB
    maxChunks    = 100
)

var allowedTypes = map[string]bool{
    "application/pdf": true,
    "application/msword": true,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
    "application/vnd.ms-excel": true,
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
    "image/png": true,
    "image/jpeg": true,
    "image/gif": true,
}

// Detect file type from magic numbers in the header of the first chunk
func detectFileType(header []byte) string {
    if len(header) >= 4 {
        if bytes.HasPrefix(header, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
            return "image/png"
        }
        if bytes.HasPrefix(header, []byte{0xFF, 0xD8}) { // JPEG magic number
            return "image/jpeg"
        }
        if bytes.HasPrefix(header, []byte{'G', 'I', 'F', '8'}) {
            return "image/gif"
        }
        if bytes.HasPrefix(header, []byte("%PDF")) {
            return "application/pdf"
        }
        if bytes.HasPrefix(header[:2], []byte{0xD0, 0xCF}) {
            return "application/msword" // DOC file magic number
        }
        if bytes.HasPrefix(header[:4], []byte{0x50, 0x4B, 0x03, 0x04}) {
            return "application/vnd.openxmlformats-officedocument.wordprocessingml.document" // DOCX
        }
        if bytes.HasPrefix(header[:2], []byte{0x09, 0x08}) {
            return "application/vnd.ms-excel" // XLS
        }
        if bytes.HasPrefix(header[:4], []byte{0x50, 0x4B, 0x03, 0x04}) {
            return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" // XLSX
        }
    }
    return "unknown"
}

// FileUploadMiddleware validates file uploads
func FileUploadMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Retrieve chunk index and total chunks from the form data
        chunkIndex, err := strconv.Atoi(c.FormValue("chunk_index"))
        if err != nil {
            log.Printf("[ERROR] Invalid chunk index: %v", err)
            return c.JSON(http.StatusBadRequest, "Invalid chunk index")
        }

        totalChunks, err := strconv.Atoi(c.FormValue("total_chunks"))
        if err != nil {
            log.Printf("[ERROR] Invalid total chunks: %v", err)
            return c.JSON(http.StatusBadRequest, "Invalid total chunks")
        }

        // Retrieve the file from the form
        file, err := c.FormFile("file")
        if err != nil {
            log.Printf("[ERROR] Failed to retrieve file: %v", err)
            return c.JSON(http.StatusBadRequest, "File is required")
        }

        // Open the file
        src, err := file.Open()
        if err != nil {
            log.Printf("[ERROR] Failed to open file: %v", err)
            return c.JSON(http.StatusInternalServerError, "Failed to open file")
        }
        defer src.Close()

        buffer := make([]byte, maxChunkSize)
        var firstChunkHeader []byte
        totalSize := int64(0) // Initialize totalSize

        log.Printf("[INFO] Starting file upload process")

        // Read file chunks
        bytesRead, err := src.Read(buffer)
        if err != nil && err != io.EOF {
            log.Printf("[ERROR] Failed to read chunk: %v", err)
            return c.JSON(http.StatusInternalServerError, "Error reading chunk")
        }

        if bytesRead > 0 {
            log.Printf("[INFO] Read chunk %d of size %d bytes", chunkIndex+1, bytesRead)
            if chunkIndex == 0 {
                // Determine file type from the first chunk
                firstChunkHeader = buffer[:bytesRead]
                if len(firstChunkHeader) > 8 {
                    firstChunkHeader = firstChunkHeader[:8] // Limit to first 8 bytes
                }
                hexHeader := hex.EncodeToString(firstChunkHeader)
                log.Printf("[INFO] First 8 bytes of chunk header (hex): 0x%s", hexHeader)

                fileType := detectFileType(firstChunkHeader)
                log.Printf("[INFO] Detected file type: %s", fileType)
                
                if !allowedTypes[fileType] {
                    log.Printf("[ERROR] Unsupported file type: %s", fileType)
                    return c.JSON(http.StatusUnsupportedMediaType, "Unsupported file type")
                }
            }

            if bytesRead > int(maxChunkSize) {
                log.Printf("[ERROR] Chunk size exceeds the maximum allowed size of %d bytes", maxChunkSize)
                return c.JSON(http.StatusRequestEntityTooLarge, "Chunk size exceeds 1 MB")
            }

            totalSize += int64(bytesRead) // Update totalSize
        }

        if chunkIndex >= totalChunks {
            log.Printf("[ERROR] Chunk index exceeds total chunks: %d/%d", chunkIndex, totalChunks)
            return c.JSON(http.StatusBadRequest, "Chunk index exceeds total chunks")
        }

        if totalChunks > maxChunks {
            log.Printf("[ERROR] Number of chunks exceeds the maximum allowed: %d", maxChunks)
            return c.JSON(http.StatusRequestEntityTooLarge, "Number of chunks exceeds 100")
        }

        log.Printf("[INFO] Total size of uploaded chunks: %d bytes", totalSize)

        // Add the total size and chunk index to the context
        c.Set("totalSize", totalSize)
        c.Set("chunkIndex", chunkIndex)
        c.Set("totalChunks", totalChunks)

        return next(c)
    }
}
