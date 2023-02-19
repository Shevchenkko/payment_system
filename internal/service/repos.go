// Package service implements application services.
package service

import (
	"context"

	// internal
	"github.com/Shevchenkko/payment_system/internal/domain"
)

///
/// REPOSITORIES
///

// Repositories contains all available repositories.
type Repositories struct {
	Users    UsersRepo
	Banks    BankAccountsRepo
	Payments PaymentsRepo
	Messages MessageLogsRepo
}

// UsersRepo - represents users repository interface.
type UsersRepo interface {
	GetUser(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, userId int) (*domain.User, error)
	CreateUser(ctx context.Context, inp *RegisterUserInput) (*domain.User, error)
	GetToken(ctx context.Context, token string) (*domain.UserToken, error)
	CreateToken(ctx context.Context, inp GenerateTokenInput) error
	DeleteToken(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, inp *ResetPasswordInput) error
}

// GenerateTokenInput represents input used to generate token.
type GenerateTokenInput struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// ResetPasswordInput - used to parameterize ResetPassword.
type ResetPasswordInput struct {
	Token    string
	Password string
}

type BankAccountsRepo interface {
	CreateBankAccount(ctx context.Context, inp *BankAccountInput, client string) (*domain.BankAccount, error)
	TopUpBankAccount(ctx context.Context, inp *TopUpBankAccountInput, balance float64) error
	CheckCreditCard(ctx context.Context, cardNumber int64) (*domain.BankAccount, error)
	GetInfoByIBAN(ctx context.Context, IBAN string) (*domain.BankAccount, error)
	ChangeCreditCardStatus(ctx context.Context, cardNumber int64, status string) (string, error)
}

type PaymentsRepo interface {
	CreatePayment(ctx context.Context, inp *PaymentInput, client *domain.BankAccount) (*domain.Payment, error)
	SentPayment(ctx context.Context, paymentId int) (string, error)
	GetPaymentByID(ctx context.Context, paymentId int) (*domain.Payment, error)
}

type MessageLogsRepo interface {
	CreateMessageLog(ctx context.Context, inp *MessageLogInput) (*domain.MessageLog, error)
}
