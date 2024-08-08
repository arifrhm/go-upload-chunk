package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var (
    AppPort    string
    UploadPath string
)

func LoadConfig() {
    // Load environment variables from .env file if present
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Read environment variables
    AppPort = os.Getenv("APP_PORT")
    if AppPort == "" {
        AppPort = "8000" // default port
    }

    UploadPath = os.Getenv("UPLOAD_PATH")
    if UploadPath == "" {
        UploadPath = "./uploads" // default upload path
    }
}
