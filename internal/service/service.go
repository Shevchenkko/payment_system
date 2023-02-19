// Package service implements application services.
package service

import (
	"context"
	"encoding/json"
	"fmt"
)

///
/// ERRORS
///

// Error - represents common struct for error
type Error struct {
	// code that will enable the API Consumers to handle type of errors
	Code string `json:"code,omitempty"`
	// message that gives the API Consumers easy-to-read explanation what went wrong and how to recover from it
	Message string `json:"message,omitempty"`
	// http status code
	Status int `json:"status,omitempty"`
}

// custom Error() method for Error
func (err *Error) Error() string {
	errData, e := json.Marshal(err)
	if e != nil {
		return fmt.Sprintf("failed to marshal error, %s", e)
	}

	return string(errData)
}

///
/// SERVICES
///

// Services contains all available services.
type Services struct {
	Users
	BankAccounts
}

// Users - represents users service interface.
type Users interface {
	RegisterUser(ctx context.Context, inp *RegisterUserInput) (RegisterUserOutput, error)
	LoginUser(ctx context.Context, inp *LoginUserInput) (LoginUserOutput, error)
	VerifyAccessToken(ctx context.Context, token string) (bool, int, string)
	SendEmail(ctx context.Context, inp *SendUserEmailInput) error
	ResetPassword(ctx context.Context, inp *ResetPasswordInput) error
}

// RegisterUserInput represents input used to register user.
type RegisterUserInput struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUserOutput - output of RegisterUser.
type RegisterUserOutput struct {
	Token    string
	UserID   int
	FullName string
	Email    string
}

// LoginUserInput represents input used to login user.
type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUserOutput - output of LoginUser.
type LoginUserOutput struct {
	Token    string
	UserID   int
	FullName string
	Email    string
}

// SendUserEmailInput represents input used to send user email.
type SendUserEmailInput struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type BankAccounts interface {
	CreateBankAccount(ctx context.Context, userId int, inp *BankAccountInput) (BankAccountOutput, error)
	BlockBankAccount(ctx context.Context, userRole string, inp *ChangeBankAccountInput) (*string, error)
	UnlockBankAccount(ctx context.Context, userRole string, inp *ChangeBankAccountInput) (*string, error)
}

type BankAccountInput struct {
	ITN         int64  `json:"itn"`
	SecretValue string `json:"secretValue"`
}

type BankAccountOutput struct {
	Client     string  `json:"client"`
	CardNumber int64   `json:"cardNumber"`
	IBAN       string  `json:"iban"`
	Balance    float64 `json:"balance"`
}

type ChangeBankAccountInput struct {
	CardNumber  int64  `json:"cardNumber"`
	SecretValue string `json:"secretValue"`
}
