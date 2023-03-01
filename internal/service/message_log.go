package service

import (
	"context"

	// internal
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

// CreateMessageLog is used to create message log.
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

// SearchLogs is used for search logs.
func (m *MessageLogsService) SearchLogs(ctx context.Context, filter *domain.Filter, client string, role string) (*SearchLogs, error) {
	if filter == nil {
		filter = new(domain.Filter)
		filter.Validate()
	}

	// search logs from db
	response, err := m.repos.Messages.SearchLogs(ctx, filter, client, role)
	if err != nil {
		return nil, err
	}

	return response, nil
}
