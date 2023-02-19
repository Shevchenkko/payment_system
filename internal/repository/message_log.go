package repository

import (
	"context"

	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/service"
	"github.com/Shevchenkko/payment_system/pkg/mysql"
)

// MessageLogsRepo - represents message log repository.
type MessageLogsRepo struct {
	*mysql.MySQL
}

// NewMessageLogsRepo - create new instance of message logs repo.
func NewMessageLogsRepo(mysql *mysql.MySQL) *MessageLogsRepo {
	return &MessageLogsRepo{mysql}
}

// CreateMessageLog - used to create message log in the database.
func (m *MessageLogsRepo) CreateMessageLog(ctx context.Context, inp *service.MessageLogInput) (*domain.MessageLog, error) {
	message := &domain.MessageLog{
		Client:  inp.Client,
		Message: inp.MessageLog,
	}
	err := m.DB.
		WithContext(ctx).
		Create(message).
		Error
	if err != nil {
		return nil, err
	}

	return message, nil
}
