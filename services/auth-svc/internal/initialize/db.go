package initialize

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("Failed to connect to database", zap.Error(err))
		return nil
	}

	return db
}

func Run() *fiber.App {
	app := fiber.New()

	return app
}
