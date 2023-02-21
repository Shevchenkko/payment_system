package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/service"
	"github.com/Shevchenkko/payment_system/pkg/mysql"
	"gorm.io/gorm"
)

// PaymentsRepo - represents payments repository.
type PaymentsRepo struct {
	*mysql.MySQL
}

// NewPaymentsRepo - create new instance of payments repo.
func NewPaymentsRepo(mysql *mysql.MySQL) *PaymentsRepo {
	return &PaymentsRepo{mysql}
}

// Search payment - used to search payment from the database.
func (p *PaymentsRepo) SearchPayments(ctx context.Context, filter *domain.Filter, client string) (*service.SearchPayments, error) {
	if filter == nil {
		filter = new(domain.Filter)
		filter.Validate()
	}

	q := p.DB.
		Table("payments").
		Offset((filter.Page - 1) * filter.List).
		Limit(filter.List).
		Order(filter.OrderString())

	var paymentOutput []service.PaymentOutput
	var response *service.SearchPayments
	if err := q.Where("from_client = ?", client).
		Find(&paymentOutput).Error; err != nil {
		return nil, &service.Error{Message: "Payments not found"}
	}

	var count int64
	q = p.DB.
		Table("payments")
	if err := q.Count(&count).Error; err != nil {
		return nil, &service.Error{Message: "Payments not found"}
	}

	response = &service.SearchPayments{
		Data: paymentOutput,
		Pagination: &domain.Pagination{
			Order: filter.OrderString(),
			Page:  filter.Page,
			List:  filter.List,
			Total: &count,
		},
	}

	return response, nil
}

// CreatePayment - used to create payment in the database.
func (p *PaymentsRepo) CreatePayment(ctx context.Context, inp *service.PaymentInput, client *domain.BankAccount) (*domain.Payment, error) {
	payment := &domain.Payment{
		FromClient:           client.Client,
		FromClientITN:        client.ITN,
		FromClientIBAN:       client.IBAN,
		FromClientCardNumber: client.CardNumber,
		Description:          inp.Description,
		ToClientIBAN:         inp.ToClientIBAN,
		ToClient:             inp.ToClient,
		OperationAmount:      inp.OperationAmount,
	}

	err := p.DB.
		WithContext(ctx).
		Create(payment).
		Error
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPaymentByID is used to get payment by id from the database.
func (p *PaymentsRepo) GetPaymentByID(ctx context.Context, paymentId int64) (*domain.Payment, error) {
	var payment domain.Payment
	err := p.DB.
		WithContext(ctx).
		Where("id = ?", paymentId).
		First(&payment).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "Payment not found"}
		}
		return nil, err
	}

	return &payment, nil
}

// SentPayment is used to sent payment.
func (p *PaymentsRepo) SentPayment(ctx context.Context, paymentId int64) (string, error) {
	status := "sent"
	err := p.DB.
		Model(domain.Payment{}).
		Where("id = ?", paymentId).
		Update("payment_status", status).
		Error
	if err != nil {
		return "", err
	}
	updatedStatus := fmt.Sprintf("Status changed to %s", status)

	return updatedStatus, err
}
