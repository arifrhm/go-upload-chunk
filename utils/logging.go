package utils

import (
    "fmt"
    "io"
    "log"
    "strings"
    "github.com/labstack/echo/v4"
    "io/ioutil"
)

const maxLogLength = 200

// CustomLogMiddleware returns a middleware function for logging request details
func CustomLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        req := c.Request()
        method := req.Method
        uri := req.RequestURI
        headers := req.Header

        // Read and limit the body size for logging
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
            body = []byte("Error reading body")
        }
        req.Body = ioutil.NopCloser(strings.NewReader(string(body))) // Restore the body for further processing

        // Prepare log entry
        logEntry := fmt.Sprintf("Method: %s, URI: %s, Headers: %v, Body: %s",
            method,
            uri,
            headers,
            truncateString(string(body), maxLogLength),
        )

        // Log to both console and file
        log.Println(logEntry)
        
        return next(c)
    }
}

// Truncate string to max length
func truncateString(s string, maxLength int) string {
    if len(s) > maxLength {
        return s[:maxLength] + "..."
    }
    return s
}

// SetLoggerOutput sets the output of the default logger to a multi-writer
func SetLoggerOutput(multiWriter io.Writer) {
    log.SetOutput(multiWriter)
}
