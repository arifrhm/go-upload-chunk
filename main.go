package main

import (
    "log"
    "os"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/arifrhm/go-upload-chunk/config"
    "github.com/arifrhm/go-upload-chunk/handlers"
    "github.com/arifrhm/go-upload-chunk/utils"
)

func main() {
    // Load the configuration
    config.LoadConfig()

    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    if _, err := os.Stat(config.UploadPath); os.IsNotExist(err) {
        // Initialize the upload directory
        utils.InitUploadDirectory()   
    }

    e.GET("/upload", handlers.UploadHandler)
    e.GET("/resume-upload", handlers.ResumeUploadHandler)
    e.POST("/upload-chunk", handlers.UploadChunkHandler)

    log.Fatal(e.Start("127.0.0.1:" + config.AppPort))
}
