package service

import (
	"context"
	"encoding/json"
	"fmt"

	// internal
	"github.com/Shevchenkko/payment_system/internal/domain"
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
	Payments
	MessageLogs
}

// Users - represents users service interface.
type Users interface {
	SearchUsers(ctx context.Context, filter *domain.Filter) (*SearchUsers, error)
	RegisterUser(ctx context.Context, inp *RegisterUserInput) (RegisterUserOutput, error)
	LoginUser(ctx context.Context, inp *LoginUserInput) (LoginUserOutput, error)
	VerifyAccessToken(ctx context.Context, token string) (bool, int, string)
	SendEmail(ctx context.Context, inp *SendUserEmailInput) error
	ResetPassword(ctx context.Context, inp *ResetPasswordInput) error
	LockUser(ctx context.Context, userId int64, userRole string) (string, error)
	UnlockUser(ctx context.Context, userId int64, userRole string) (string, error)
}

// SearchUsers represents user info.
type SearchUsers struct {
	Data       []User             `json:"data"`
	Pagination *domain.Pagination `json:"pagination"`
}

type User struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	FullName string `json:"fullName"`
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

// BankAccounts - represents bank accounts service interface.
type BankAccounts interface {
	SearchBankAccounts(ctx context.Context, filter *domain.Filter, client string, role string) (*SearchBankAccounts, error)
	CreateBankAccount(ctx context.Context, userId int, inp *BankAccountInput) (BankAccountOutput, error)
	TopUpBankAccount(ctx context.Context, userId int, inp *TopUpBankAccountInput) (BankAccountOutput, error)
	BlockBankAccount(ctx context.Context, client string, userRole string, inp *ChangeBankAccountInput) (string, error)
	UnlockBankAccount(ctx context.Context, client string, userRole string, inp *ChangeBankAccountInput) (string, error)
}

// BankAccountInput represents input used to bank account.
type BankAccountInput struct {
	ITN         int64  `json:"itn"`
	SecretValue string `json:"secretValue"`
}

// SearchBankAccounts represents bank account info.
type SearchBankAccounts struct {
	Data       []BankAccountOutput `json:"data"`
	Pagination *domain.Pagination  `json:"pagination"`
}

type BankAccountOutput struct {
	ID         int     `json:"id"`
	Status     string  `json:"status"`
	Client     string  `json:"client"`
	CardNumber int64   `json:"cardNumber"`
	IBAN       string  `json:"iban"`
	Balance    float64 `json:"balance"`
}

// TopUpBankAccountInput represents input used to top up bank account.
type TopUpBankAccountInput struct {
	CardNumber      int64   `json:"cardNumber"`
	OperationAmount float64 `json:"operationAmount"`
}

// ChangeBankAccountInput represents input used to change bank account.
type ChangeBankAccountInput struct {
	CardNumber  int64  `json:"cardNumber"`
	SecretValue string `json:"secretValue"`
}

// Payments - represents payments service interface.
type Payments interface {
	SearchPayments(ctx context.Context, filter *domain.Filter, client string) (*SearchPayments, error)
	CreatePayment(ctx context.Context, userId int, inp *PaymentInput) (*PaymentOutput, error)
	SentPayment(ctx context.Context, paymentId int64, secretValue string, cardBalance float64) (string, error)
}

// PaymentInput represents input used to payment.
type PaymentInput struct {
	FromClientIBAN  string  `json:"fromClientIban"`
	Description     string  `json:"description"`
	ToClientIBAN    string  `json:"toClientIban"`
	ToClient        string  `json:"toClient"`
	OperationAmount float64 `json:"operationAmount"`
}

// SearchPayments represents payments info.
type SearchPayments struct {
	Data       []PaymentOutput    `json:"data"`
	Pagination *domain.Pagination `json:"pagination"`
}

type PaymentOutput struct {
	ID                   int64   `json:"id"`
	PaymentStatus        string  `json:"paymentStatus"`
	FromClient           string  `json:"fromClient"`
	FromClientITN        int64   `json:"fromClientItn"`
	FromClientIBAN       string  `json:"fromClientIban"`
	FromClientCardNumber int64   `json:"fromClientCardNumber"`
	Description          string  `json:"description"`
	ToClientIBAN         string  `json:"toClientIban"`
	ToClient             string  `json:"toClient"`
	OperationAmount      float64 `json:"operationAmount"`
}

// MessageLogs - represents message logs service interface.
type MessageLogs interface {
	CreateMessageLog(ctx context.Context, userId int, inp *MessageLogInput) (*domain.MessageLog, error)
	SearchLogs(ctx context.Context, filter *domain.Filter, client string, role string) (*SearchLogs, error)
}

// MessageLogInput represents input used to message logs.
type MessageLogInput struct {
	Client     string `json:"client"`
	MessageLog string `json:"messageLog"`
}

// SearchLogs represents logs info.
type SearchLogs struct {
	Data       []domain.MessageLog `json:"data"`
	Pagination *domain.Pagination  `json:"pagination"`
}
