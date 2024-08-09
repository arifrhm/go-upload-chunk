package utils

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "strings"
    "time"
    "github.com/labstack/echo/v4"
    "io/ioutil"
)

const maxLogLength = 200 // Define maxLogLength

// LogEntry represents the structure of the log entry
type LogEntry struct {
    Time         string `json:"time"`
    ID           string `json:"id"`
    RemoteIP     string `json:"remote_ip"`
    Host         string `json:"host"`
    Method       string `json:"method"`
    URI          string `json:"uri"`
    UserAgent    string `json:"user_agent"`
    Status       int    `json:"status"`
    Error        string `json:"error"`
    Latency      int64  `json:"latency"`
    LatencyHuman string `json:"latency_human"`
    BytesIn      int64  `json:"bytes_in"`
    BytesOut     int64  `json:"bytes_out"`
    Success      bool   `json:"success,omitempty"`
    Data         string `json:"data,omitempty"`
    Params       string `json:"params,omitempty"`
    RequestID    string `json:"request_id,omitempty"`
}

// CustomLogMiddleware returns a middleware function for logging request details in JSON format
func CustomLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        req := c.Request()
        res := c.Response()

        // Read and limit the body size for logging
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
            body = []byte("Error reading body")
        }
        req.Body = ioutil.NopCloser(strings.NewReader(string(body))) // Restore the body for further processing

        // Process the request
        err = next(c)

        // Calculate latency
        stop := time.Now()
        latency := stop.Sub(start)

        // Prepare log entry
        logEntry := LogEntry{
            Time:         stop.Format(time.RFC3339Nano),
            ID:           c.Response().Header().Get(echo.HeaderXRequestID),
            RemoteIP:     c.RealIP(),
            Host:         req.Host,
            Method:       req.Method,
            URI:          req.RequestURI,
            UserAgent:    req.UserAgent(),
            Status:       res.Status,
            Error:        "",
            Latency:      latency.Microseconds(),
            LatencyHuman: latency.String(),
            BytesIn:      req.ContentLength,
            BytesOut:     res.Size,
            Success:      res.Status < 400,
            Data:         truncateString(string(body), maxLogLength),
            Params:       fmt.Sprintf("%v", c.QueryParams()),
            RequestID:    c.Response().Header().Get(echo.HeaderXRequestID),
        }

        if err != nil {
            logEntry.Error = err.Error()
        }

        logJSON, _ := json.Marshal(logEntry)

        // Log to both console and file
        log.Println(string(logJSON))

        return err
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
