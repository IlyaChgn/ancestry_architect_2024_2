package repository

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	ctx := context.Background()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	manager := NewSessionManager(db)

	userID := uint(123)
	sessionID := uuid.NewString()

	// Устанавливаем ожидание, что будет вызван метод Set с указанными значениями
	mock.ExpectSet(sessionID, userID, 3*24*time.Hour).SetVal("OK")

	err := manager.CreateSession(ctx, sessionID, userID)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRemoveExistingSession(t *testing.T) {
	ctx := context.Background()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	manager := NewSessionManager(db)

	userID := uint(123)
	sessionID := uuid.NewString()

	mock.ExpectSet(sessionID, userID, 3*24*time.Hour).SetVal("OK")
	manager.CreateSession(ctx, sessionID, userID)

	mock.ExpectGet(sessionID).SetVal("123")
	mock.ExpectDel(sessionID).SetVal(1)
	err := manager.RemoveSession(ctx, sessionID)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRemoveNonExistingSession(t *testing.T) {
	ctx := context.Background()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	manager := NewSessionManager(db)

	sessionID := uuid.NewString()

	mock.ExpectGet(sessionID).SetVal("0")
	err := manager.RemoveSession(ctx, sessionID)

	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
