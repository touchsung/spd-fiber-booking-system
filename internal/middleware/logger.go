package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelError LogLevel = "ERROR"
)

type RequestLog struct {
	Level      LogLevel    `json:"level"`
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id"`
	Method     string      `json:"method"`
	Path       string      `json:"path"`
	IP         string      `json:"ip"`
	StatusCode int         `json:"status_code"`
	Duration   string      `json:"duration"`
	Body       interface{} `json:"body,omitempty"`
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())
		c.Locals("requestID", requestID)

		// Parse body if exists
		var body interface{}
		if len(c.Body()) > 0 {
			if err := json.Unmarshal(c.Body(), &body); err == nil {
				c.Locals("requestBody", body)
			}
		}

		err := c.Next()

		// Create log entry
		logEntry := RequestLog{
			Level:      LogLevelInfo,
			Timestamp:  time.Now(),
			RequestID:  requestID,
			Method:     c.Method(),
			Path:       c.Path(),
			IP:         c.IP(),
			StatusCode: c.Response().StatusCode(),
			Duration:   time.Since(start).String(),
			Body:       c.Locals("requestBody"),
		}

		// If there's an error, change log level
		if err != nil || c.Response().StatusCode() >= 400 {
			logEntry.Level = LogLevelError
		}

		// Convert to JSON and print
		logJSON, _ := json.Marshal(logEntry)
		fmt.Println(string(logJSON))

		return err
	}
}
