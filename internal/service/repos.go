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
	Users UsersRepo
	Banks BankAccountsRepo
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
	CheckCreditCard(ctx context.Context, cardNumber int64) (*domain.BankAccount, error)
	ChangeCreditCardStatus(ctx context.Context, cardNumber int64, status string) (string, error)
}
