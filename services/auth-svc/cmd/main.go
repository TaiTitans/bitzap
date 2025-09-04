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
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/taititans/bitzap/auth-svc/internal/config"
	"github.com/taititans/bitzap/auth-svc/internal/controller/http"
	"github.com/taititans/bitzap/auth-svc/internal/controller/http/auth"
	"github.com/taititans/bitzap/auth-svc/internal/controller/http/email"
	repository_impl "github.com/taititans/bitzap/auth-svc/internal/domain/repository/repository_impl"
	"github.com/taititans/bitzap/auth-svc/internal/initialize"
	"github.com/taititans/bitzap/auth-svc/internal/logic"
	"github.com/taititans/bitzap/auth-svc/internal/middleware"
	"github.com/taititans/bitzap/auth-svc/internal/service"
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
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration from environment variables
	cfg := config.LoadConfig()

	loggerConfig := initialize.LoggerConfig{
		Path:   "./log/",
		File:   "app.log",
		Level:  cfg.Server.LogLevel,
		Stdout: true,
		StSkip: 1,
	}

	logger := initialize.InitLogger(loggerConfig)
	defer logger.Sync() // Flush buffer

	// logger wrapper
	appLogger := util.NewZapLogger(logger)
	appLogger.Info("Starting Auth Service")

	// Database configuration from environment
	dbConfig := initialize.DatabaseConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.Database.ConnMaxLifetime) * time.Hour,
	}

	// Initialize database
	db := initialize.InitDatabase(dbConfig)
	defer initialize.CloseDatabase(db)
	appLogger.Info("POSTGRES database connected successfully: auth-service")

	// Run migrations
	// if err := initialize.AutoMigrate(db); err != nil {
	// 	appLogger.Error("Failed to run database migrations", util.Error(err))
	// 	log.Fatal(err)
	// }
	appLogger.Info("Database migrations completed")

	// Initialize repositories
	userRepo := repository_impl.NewUserRepository(db)
	userRoleRepo := repository_impl.NewUserRoleRepository(db)
	userPermissionRepo := repository_impl.NewUserPermissionRepository(db)
	userActivityLogRepo := repository_impl.NewUserActivityLogRepository(db)

	// Redis configuration from environment
	redisConfig := initialize.RedisConfig{
		Address:      cfg.Redis.Default.Address,
		Password:     cfg.Redis.Default.Password,
		DB:           cfg.Redis.Default.DB,
		DialTimeout:  cfg.Redis.Default.DialTimeout,
		ReadTimeout:  cfg.Redis.Default.ReadTimeout,
		WriteTimeout: cfg.Redis.Default.WriteTimeout,
		MaxActive:    cfg.Redis.Default.MaxActive,
	}

	redisClient := initialize.InitRedis(redisConfig)
	defer initialize.CloseRedis(redisClient)
	appLogger.Info("Redis connected successfully")

	// Initialize email service with config from environment
	emailService := service.NewEmailService(cfg.Email, redisClient, appLogger)

	// Initialize business logic
	authLogic := logic.NewAuthLogic(userRepo, userRoleRepo, userPermissionRepo, userActivityLogRepo, emailService, appLogger)

	// Initialize services
	authService := service.NewAuthService(authLogic)

	// Initialize controllers
	authController := auth.NewAuthController(authService, appLogger)
	emailController := email.NewEmailController(emailService, appLogger)

	// Fiber app
	app := fiber.New()

	// Middleware
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggerMiddleware(appLogger))
	app.Use(middleware.ErrorLoggerMiddleware(appLogger))

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Setup auth routes
	http.SetupAuthRoutes(app, authController, emailController)

	// Ping route
	app.Get("/ping", func(c *fiber.Ctx) error {
		appLogger.Info("Ping request received")
		return c.JSON(fiber.Map{"message": "pong"})
	})

	// Log Swagger URL - sử dụng log.Printf thay vì appLogger
	swaggerURL := "http://localhost:" + cfg.Server.Port + "/swagger/"
	log.Printf("📚 Swagger documentation available at: %s", swaggerURL)

	appLogger.Info("Server starting on :" + cfg.Server.Port)
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
