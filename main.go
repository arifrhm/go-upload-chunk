package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/arifrhm/go-upload-chunk/handlers"
    "github.com/arifrhm/go-upload-chunk/config"
    "github.com/arifrhm/go-upload-chunk/utils"
    "io"
    "os"
)

func main() {
    // Create log file
    logFile, err := os.OpenFile(config.LogPath+"/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        panic("Could not open log file: " + err.Error())
    }
    // Ensure the file is closed properly
    defer logFile.Close()

    // Create a MultiWriter to write to both stdout and the log file
    multiWriter := io.MultiWriter(os.Stdout, logFile)

    // Set the logger output
    utils.SetLoggerOutput(multiWriter)

    // Create Echo instance
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Use custom log middleware
    e.Use(utils.CustomLogMiddleware)

    // Register your routes
    handlers.RegisterRoutes(e)

    // Start the server
    e.Logger.Fatal(e.Start("127.0.0.1:" + config.AppPort))
}
