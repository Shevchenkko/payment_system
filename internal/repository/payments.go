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
func (p *PaymentsRepo) GetPaymentByID(ctx context.Context, paymentId int) (*domain.Payment, error) {
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
func (p *PaymentsRepo) SentPayment(ctx context.Context, paymentId int) (string, error) {
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
