// @title           Auth Service API
// @version         1.0
// @description     Authentication service API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/taititans/bitzap/auth-svc/internal/initialize"
	"github.com/taititans/bitzap/auth-svc/internal/middleware"
	"github.com/taititans/bitzap/auth-svc/internal/util"

	_ "github.com/taititans/bitzap/auth-svc/docs" // swagger docs
)

// @Summary     Test Redis connection
// @Description Test Redis operations including set, get, and increment counter
// @Tags        redis
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]interface{} "Success response with test data and counter"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /redis/test [get]
func main() {
	// Khởi tạo logger
	loggerConfig := initialize.LoggerConfig{
		Path:   "./log/",
		File:   "app.log",
		Level:  "info",
		Stdout: true,
		StSkip: 1,
	}

	logger := initialize.InitLogger(loggerConfig)
	defer logger.Sync() // Flush buffer

	// Tạo logger wrapper
	appLogger := util.NewZapLogger(logger)
	appLogger.Info("Starting Auth Service")

	// Redis
	redisConfig := initialize.RedisConfig{
		Address:      "127.0.0.1:6379",
		Password:     "redispass",
		DB:           0,
		DialTimeout:  "30s",
		ReadTimeout:  "30s",
		WriteTimeout: "30s",
		MaxActive:    100,
	}

	redisClient := initialize.InitRedis(redisConfig)
	defer initialize.CloseRedis(redisClient)
	appLogger.Info("Redis connected successfully")

	// Redis helper
	redisHelper := util.NewRedisHelper(redisClient)

	// Fiber app
	app := fiber.New()

	// Middleware
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggerMiddleware(appLogger))
	app.Use(middleware.ErrorLoggerMiddleware(appLogger))

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Demo route test Redis
	app.Get("/redis/test", func(c *fiber.Ctx) error {
		ctx := context.Background()

		// Test set/get
		testData := map[string]interface{}{
			"message":   "Hello Redis!",
			"timestamp": time.Now().Unix(),
		}

		err := redisHelper.Set(ctx, "test:key", testData, 5*time.Minute)
		if err != nil {
			appLogger.Error("Failed to set Redis key", util.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": "Failed to set Redis key"})
		}

		// Test get
		var retrievedData map[string]interface{}
		err = redisHelper.Get(ctx, "test:key", &retrievedData)
		if err != nil {
			appLogger.Error("Failed to get Redis key", util.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": "Failed to get Redis key"})
		}

		// Test counter
		counter, err := redisHelper.Incr(ctx, "test:counter")
		if err != nil {
			appLogger.Error("Failed to increment counter", util.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": "Failed to increment counter"})
		}

		appLogger.Info("Redis test completed successfully",
			util.Int64("counter", counter),
			util.String("data", "retrieved successfully"),
		)

		return c.JSON(fiber.Map{
			"message": "Redis test successful",
			"data":    retrievedData,
			"counter": counter,
		})
	})

	// Ping route
	app.Get("/ping", func(c *fiber.Ctx) error {
		appLogger.Info("Ping request received")
		return c.JSON(fiber.Map{"message": "pong"})
	})

	appLogger.Info("Server starting on :8080")
	log.Fatal(app.Listen(":8080"))
}
