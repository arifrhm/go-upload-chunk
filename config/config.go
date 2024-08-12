package config

import (
    "github.com/joho/godotenv"
    "os"
    "strings"
)

var (
    AppPort       string
    UploadPath    string
    LogPath       string
    AllowedOrigins []string
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
    
    // Load allowed origins from environment variable
    AllowedOrigins = strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ",")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
