package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"

	logging "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/utils"
)

const (
	errSessionNotExists = "session doesn`t exist"
	SessionDuration     = 3 * 24 * time.Hour
)

type SessionManager struct {
	client *redis.Client
}

func NewSessionManager(client *redis.Client) *SessionManager {
	return &SessionManager{
		client: client,
	}
}

func (manager *SessionManager) CreateSession(ctx context.Context, sessionID string, userID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("redis", logging.GetFunctionName()))

	err := manager.client.Set(context.Background(), sessionID, userID, SessionDuration).Err()
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while setting user session, %w", err))

		return err
	}

	return nil
}

func (manager *SessionManager) RemoveSession(ctx context.Context, sessionID string) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("redis", logging.GetFunctionName()))

	if _, exists := manager.GetUserID(ctx, sessionID); !exists {
		err := fmt.Errorf("something went wrong while removing session: %s", errSessionNotExists)
		logging.LogError(logger, err)

		return err
	}

	_, err := manager.client.Del(context.Background(), sessionID).Result()
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while removing user session, %w", err))

		return err
	}

	return nil
}

func (manager *SessionManager) GetUserID(ctx context.Context, sessionID string) (uint, bool) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("redis", logging.GetFunctionName()))

	id, _ := manager.client.Get(context.Background(), sessionID).Result()

	userID, err := strconv.Atoi(id)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("error occured while parsing userID to int, %w", err))

		return 0, false
	}

	return uint(userID), userID != 0
}
