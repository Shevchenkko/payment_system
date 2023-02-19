package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// PaymentsService - represents payments service.
type PaymentsService struct {
	repos Repositories
}

// NewPaymentService - creates instance of new payment service.
func NewPaymentService(repos Repositories) *PaymentsService {
	return &PaymentsService{repos}
}

// CreatePayment is used for creating payment.
func (p *PaymentsService) CreatePayment(ctx context.Context, userId int, inp *PaymentInput) (*PaymentOutput, error) {
	client, err := p.repos.Banks.GetInfoByIBAN(ctx, inp.FromClientIBAN)
	if err != nil {
		return nil, err
	}

	// create payment in db
	payment, err := p.repos.Payments.CreatePayment(ctx, inp, client)
	if err != nil {
		return nil, err
	}

	return &PaymentOutput{
		PaymentID:            payment.ID,
		PaymentStatus:        payment.PaymentStatus,
		FromClient:           payment.FromClient,
		FromClientITN:        payment.FromClientITN,
		FromClientIBAN:       payment.FromClientIBAN,
		FromClientCardNumber: payment.FromClientCardNumber,
		Description:          payment.Description,
		ToClientIBAN:         payment.ToClientIBAN,
		ToClient:             payment.ToClient,
		OperationAmount:      payment.OperationAmount,
	}, nil
}

// SentPayment is used for senting payment.
func (p *PaymentsService) SentPayment(ctx context.Context, paymentId int, secretValue string, cardBalance float64) (string, error) {
	// check payment
	payment, err := p.repos.Payments.GetPaymentByID(ctx, paymentId)
	if err != nil {
		return "", err
	}

	// get secret value
	bakn, err := p.repos.Banks.GetInfoByIBAN(ctx, payment.FromClientIBAN)
	if err != nil {
		return "", err
	}

	// check secret value
	err = bcrypt.CompareHashAndPassword([]byte(bakn.SecretValue), []byte(secretValue))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", &Error{Message: "Wrong secret value"}
		}
		return "", err
	}

	// sent payment
	status, err := p.repos.Payments.SentPayment(ctx, paymentId)
	if err != nil {
		return "", err
	}

	// update balance
	err = p.repos.Banks.TopUpBankAccount(ctx, &TopUpBankAccountInput{
		CardNumber: bakn.CardNumber}, cardBalance)
	if err != nil {
		return "", err
	}

	return status, nil
}