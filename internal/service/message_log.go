package service

import (
	"context"

	"github.com/Shevchenkko/payment_system/internal/domain"
)

// MessageLogs - represents message log service.
type MessageLogsService struct {
	repos Repositories
}

// NewMessageLogsService - creates instance of new message logs service.
func NewMessageLogsService(repos Repositories) *MessageLogsService {
	return &MessageLogsService{repos}
}

// CreateMessageLog - used to create message log.
func (m *MessageLogsService) CreateMessageLog(ctx context.Context, userId int, inp *MessageLogInput) (*domain.MessageLog, error) {
	client, err := m.repos.Users.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	// create message log in db
	message, err := m.repos.Messages.CreateMessageLog(ctx, &MessageLogInput{
		Client:     client.FullName,
		MessageLog: inp.MessageLog,
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}
