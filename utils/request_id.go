// utils/request_id.go

package utils

import (
    "github.com/labstack/echo/v4"
    "github.com/google/uuid"
)

// RequestIDMiddleware generates a unique request ID for each request and attaches it to the response headers.
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Generate a unique request ID
        requestID := uuid.New().String()

        // Add the request ID to the request context
        c.Set("request_id", requestID)

        // Add the request ID to the response headers
        c.Response().Header().Set("X-Request-ID", requestID)

        return next(c)
    }
}
