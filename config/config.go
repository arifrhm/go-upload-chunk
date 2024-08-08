package config

import (
    "github.com/joho/godotenv"
    "os"
)

var (
    AppPort    string
    UploadPath string
    LogPath    string
)

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    // Initialize configuration variables
    AppPort = getEnv("APP_PORT", "8001")
    UploadPath = getEnv("UPLOAD_PATH", "./uploads")
    LogPath = getEnv("LOG_PATH", "./logs")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
