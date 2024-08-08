package utils

import (
    "os"
    "log"
    "github.com/arifrhm/go-upload-chunk/config"
)

func InitUploadDirectory() {
    if _, err := os.Stat(config.UploadPath); os.IsNotExist(err) {
        err := os.Mkdir(config.UploadPath, os.ModePerm)
        if err != nil {
            log.Fatalf("Failed to create upload directory: %v", err)
        }
    }
}
