package middleware

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// LoggerMiddleware create middleware for log HTTP requests
func LoggerMiddleware(logger util.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request
		logger.Info("HTTP Request",
			util.Method(c.Method()),
			util.Path(c.Path()),
			util.StatusCode(c.Response().StatusCode()),
			util.Duration(duration),
			util.String("ip", c.IP()),
			util.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}

// ErrorLoggerMiddleware create middleware for log errors
func ErrorLoggerMiddleware(logger util.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			logger.Error("HTTP Error",
				util.Method(c.Method()),
				util.Path(c.Path()),
				util.StatusCode(c.Response().StatusCode()),
				util.Error(err),
				util.String("ip", c.IP()),
			)
		}

		return err
	}
}

// RequestIDMiddleware create middleware with request ID
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Set("X-Request-ID", requestID)
		}

		// Add request ID to context
		c.Locals("request_id", requestID)

		return c.Next()
	}
}

// generateRequestID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
