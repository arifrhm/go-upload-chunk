package main

import (
	"io"
	"os"

	"github.com/arifrhm/go-upload-chunk/config"
	_ "github.com/arifrhm/go-upload-chunk/docs" // Import generated docs
	"github.com/arifrhm/go-upload-chunk/handlers"
	"github.com/arifrhm/go-upload-chunk/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	e.Use(utils.RequestIDMiddleware)

	// CORS middleware using config package
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// Swagger UI route
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Register your routes
	handlers.RegisterRoutes(e)
	// Start the server
	e.Logger.Fatal(e.Start(config.AppHost + ":" + config.AppPort))
}
