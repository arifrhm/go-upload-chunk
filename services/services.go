package services

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "github.com/arifrhm/go-upload-chunk/config"
)

var tempChunkSuffix = ".part"

// CheckAndAssembleFile checks if all chunks are uploaded and assembles the final file if true.
func CheckAndAssembleFile(fileName string, totalChunks int) error {
    chunks, err := filepath.Glob(filepath.Join(config.UploadPath, fmt.Sprintf("%s%s*", fileName, tempChunkSuffix)))
    if err != nil {
        return fmt.Errorf("failed to list chunk files: %v", err)
    }

    // Check if all chunks are present
    if len(chunks) != totalChunks {
        return fmt.Errorf("not all chunks are uploaded yet")
    }

    finalPath := filepath.Join(config.UploadPath, fileName)
    finalFile, err := os.Create(finalPath)
    if err != nil {
        return fmt.Errorf("failed to create final file: %v", err)
    }
    defer finalFile.Close()

    // Assemble chunks into the final file
    for i := 0; i < totalChunks; i++ {
        chunkPath := filepath.Join(config.UploadPath, fmt.Sprintf("%s%s%d", fileName, tempChunkSuffix, i))
        chunkFile, err := os.Open(chunkPath)
        if err != nil {
            return fmt.Errorf("failed to open chunk file %d: %v", i, err)
        }
        _, err = io.Copy(finalFile, chunkFile)
        chunkFile.Close()
        if err != nil {
            return fmt.Errorf("failed to copy chunk file %d to final file: %v", i, err)
        }
        os.Remove(chunkPath)
    }
    return nil
}
