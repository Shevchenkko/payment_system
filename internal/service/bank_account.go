package service

import (
	"context"
	"errors"

	"github.com/Shevchenkko/payment_system/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// BankAccountsService - represents bank accounts service.
type BankAccountsService struct {
	repos Repositories
}

// NewBankAccountService - creates instance of new bank account service.
func NewBankAccountService(repos Repositories) *BankAccountsService {
	return &BankAccountsService{repos}
}

// SearchBankAccount is used for search bank account.
func (b *BankAccountsService) SearchBankAccounts(ctx context.Context, filter *domain.Filter, client string, role string) (*SearchBankAccounts, error) {
	if filter == nil {
		filter = new(domain.Filter)
		filter.Validate()
	}

	// search bank accounts from db
	response, err := b.repos.Banks.SearchBankAccounts(ctx, filter, client, role)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateBankAccount is used for creating bank account.
func (b *BankAccountsService) CreateBankAccount(ctx context.Context, userId int, inp *BankAccountInput) (BankAccountOutput, error) {
	client, err := b.repos.Users.GetUserByID(ctx, userId)
	if err != nil {
		return BankAccountOutput{}, err
	}

	// create bank account in db
	account, err := b.repos.Banks.CreateBankAccount(ctx, inp, client.FullName)
	if err != nil {
		return BankAccountOutput{}, err
	}

	return BankAccountOutput{
		ID:         account.ID,
		Client:     client.FullName,
		CardNumber: account.CardNumber,
		IBAN:       account.IBAN,
	}, nil
}

// TopUpBankAccount is used for top up bank account.
func (b *BankAccountsService) TopUpBankAccount(ctx context.Context, userId int, inp *TopUpBankAccountInput) (BankAccountOutput, error) {
	client, err := b.repos.Users.GetUserByID(ctx, userId)
	if err != nil {
		return BankAccountOutput{}, err
	}

	// get card
	card, err := b.repos.Banks.CheckCreditCard(ctx, inp.CardNumber)
	if err != nil {
		return BankAccountOutput{}, err
	}

	cardBalance := card.Balance + inp.OperationAmount

	// top up bank account in db
	err = b.repos.Banks.TopUpBankAccount(ctx, inp, cardBalance)
	if err != nil {
		return BankAccountOutput{}, err
	}

	return BankAccountOutput{
		Client:     client.FullName,
		CardNumber: card.CardNumber,
		IBAN:       card.IBAN,
		Balance:    cardBalance,
	}, nil
}

// BlockBankAccount is used for blocing bank account.
func (b *BankAccountsService) BlockBankAccount(ctx context.Context, client string, userRole string, inp *ChangeBankAccountInput) (string, error) {
	status, err := b.repos.Banks.CheckCreditCard(ctx, inp.CardNumber)
	if err != nil {
		return "", err
	}
	var accountChange string

	// check user role
	if userRole == "admin" {
		if status.Status == "ACTIVE" {
			accountChange, err = b.repos.Banks.ChangeCreditCardStatus(ctx, inp.CardNumber, "LOCK")
			if err != nil {
				return "", err
			}
		} else {
			accountChange = "The account has already been blocked"
		}
	} else {
		if client == status.Client {
			// check secret value
			err = bcrypt.CompareHashAndPassword([]byte(status.SecretValue), []byte(inp.SecretValue))
			if err != nil {
				if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
					return "", &Error{Message: "Wrong secret value"}
				}
				return "", err
			}
			if status.Status == "ACTIVE" {
				accountChange, err = b.repos.Banks.ChangeCreditCardStatus(ctx, inp.CardNumber, "LOCK")
				if err != nil {
					return "", err
				}
			} else {
				accountChange = "The account has already been blocked"
			}
		} else {
			accountChange = "This is not your bank account!"
		}
	}
	return accountChange, nil
}

// UnlockBankAccount is used for unlocing bank account.
func (b *BankAccountsService) UnlockBankAccount(ctx context.Context, client string, userRole string, inp *ChangeBankAccountInput) (string, error) {
	status, err := b.repos.Banks.CheckCreditCard(ctx, inp.CardNumber)
	if err != nil {
		return "", err
	}
	var accountChange string

	// check user role
	if userRole == "admin" {
		if status.Status == "LOCK" {
			accountChange, err = b.repos.Banks.ChangeCreditCardStatus(ctx, inp.CardNumber, "ACTIVE")
			if err != nil {
				return "", err
			}
		} else {
			accountChange = "The account has already been active"
		}
	} else {
		if client == status.Client {
			// check secret value
			err = bcrypt.CompareHashAndPassword([]byte(status.SecretValue), []byte(inp.SecretValue))
			if err != nil {
				if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
					return "", &Error{Message: "Wrong secret value"}
				}
				return "", err
			}
			if status.Status == "LOCK" {
				accountChange, err = b.repos.Banks.ChangeCreditCardStatus(ctx, inp.CardNumber, "ACTIVE")
				if err != nil {
					return "", err
				}
			} else {
				accountChange = "The account has already been active"
			}
		} else {
			accountChange = "This is not your bank account!"
		}
	}
	return accountChange, nil
}
