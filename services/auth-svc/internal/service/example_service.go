package service

import (
	"context"
	"time"

	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// ExampleService service with logger
type ExampleService struct {
	logger util.Logger
}

// NewExampleService create ExampleService
func NewExampleService(logger util.Logger) *ExampleService {
	return &ExampleService{
		logger: logger,
	}
}

// ProcessData with logging
func (s *ExampleService) ProcessData(ctx context.Context, data string) (string, error) {
	s.logger.Info("Processing data",
		util.String("data", data),
		util.String("operation", "process"),
	)

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Simulate some business logic
	if data == "error" {
		s.logger.Error("Data processing failed",
			util.String("data", data),
			util.String("reason", "invalid_data"),
		)
		return "", util.NewError("invalid data")
	}

	result := "processed_" + data
	s.logger.Info("Data processed successfully",
		util.String("input", data),
		util.String("output", result),
	)

	return result, nil
}

// GetUserInfo
func (s *ExampleService) GetUserInfo(ctx context.Context, userID string) (map[string]interface{}, error) {
	s.logger.Info("Getting user info",
		util.String("user_id", userID),
	)

	// Simulate database query
	time.Sleep(50 * time.Millisecond)

	userInfo := map[string]interface{}{
		"id":    userID,
		"name":  "John Doe",
		"email": "john@example.com",
	}

	s.logger.Info("User info retrieved",
		util.String("user_id", userID),
		util.String("name", userInfo["name"].(string)),
	)

	return userInfo, nil
}

// UpdateUser
func (s *ExampleService) UpdateUser(ctx context.Context, userID string, data map[string]interface{}) error {
	s.logger.Info("Updating user",
		util.String("user_id", userID),
		util.Any("data", data),
	)

	// Simulate database update
	time.Sleep(200 * time.Millisecond)

	s.logger.Info("User updated successfully",
		util.String("user_id", userID),
	)

	return nil
}
