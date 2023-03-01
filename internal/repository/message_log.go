package repository

import (
	"context"

	// external
	"github.com/Shevchenkko/payment_system/pkg/mysql"

	// internal
	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/service"
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

// Search logs - used to search log from the database.
func (m *MessageLogsRepo) SearchLogs(ctx context.Context, filter *domain.Filter, client string, role string) (*service.SearchLogs, error) {
	if filter == nil {
		filter = new(domain.Filter)
		filter.Validate()
	}

	q := m.DB.
		Table("message_logs").
		Offset((filter.Page - 1) * filter.List).
		Limit(filter.List).
		Order(filter.OrderString())

	var logOutput []domain.MessageLog
	var response *service.SearchLogs
	if role == "admin" {
		if err := q.Find(&logOutput).Error; err != nil {
			return nil, &service.Error{Message: "Logs not found"}
		}
	} else {
		if err := q.Where("client = ?", client).
			Find(&logOutput).Error; err != nil {
			return nil, &service.Error{Message: "Logs not found"}
		}
	}

	var count int64
	q = m.DB.
		Table("message_logs")
	if err := q.Count(&count).Error; err != nil {
		return nil, &service.Error{Message: "Logs not found"}
	}

	response = &service.SearchLogs{
		Data: logOutput,
		Pagination: &domain.Pagination{
			Order: filter.OrderString(),
			Page:  filter.Page,
			List:  filter.List,
			Total: &count,
		},
	}

	return response, nil
}
