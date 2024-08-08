package services

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "mime/multipart"
    "github.com/arifrhm/go-upload-chunk/config"
    "github.com/arifrhm/go-upload-chunk/utils"
)

var tempChunkSuffix = ".part"

// CheckAndAssembleFile checks if all chunks are uploaded and assembles the final file.
func CheckAndAssembleFile(fileName string, totalChunks int) error {
    chunks, err := filepath.Glob(filepath.Join(config.UploadPath, fmt.Sprintf("%s%s*", fileName, tempChunkSuffix)))
    if err != nil {
        return err
    }

    if len(chunks) == totalChunks {
        finalPath := filepath.Join(config.UploadPath, fileName)
        finalFile, err := os.Create(finalPath)
        if err != nil {
            return err
        }
        defer finalFile.Close()

        for i := 0; i < totalChunks; i++ {
            chunkPath := filepath.Join(config.UploadPath, fmt.Sprintf("%s%s%d", fileName, tempChunkSuffix, i))
            chunkFile, err := os.Open(chunkPath)
            if err != nil {
                return err
            }
            _, err = io.Copy(finalFile, chunkFile)
            chunkFile.Close()
            if err != nil {
                return err
            }
            os.Remove(chunkPath)
        }
        return nil
    }
    return fmt.Errorf("not all chunks are uploaded yet")
}

// HandleFileUpload handles the file upload and validates the file type only for the first chunk.
func HandleFileUpload(file multipart.File, fileHeader *multipart.FileHeader, fileName string, chunkIndex, totalChunks int) error {
    if chunkIndex == 0 {
        if err := utils.ValidateFileType(fileHeader); err != nil {
            return err
        }
    }

    if err := utils.ValidateFileSize(file); err != nil {
        return err
    }

    chunkPath := filepath.Join(config.UploadPath, fmt.Sprintf("%s%s%d", fileName, tempChunkSuffix, chunkIndex))

    dst, err := os.Create(chunkPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        return err
    }

    // Check and assemble file only if it's the last chunk
    if chunkIndex == totalChunks-1 {
        return CheckAndAssembleFile(fileName, totalChunks)
    }

    return nil
}
