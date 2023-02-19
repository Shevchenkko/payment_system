package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/service"
	"github.com/Shevchenkko/payment_system/pkg/mysql"
	"github.com/Shevchenkko/payment_system/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// BankAccountsRepo - represents bank account repository.
type BankAccountsRepo struct {
	*mysql.MySQL
}

// NewBankAccountsRepo - create new instance of bank accounts repo.
func NewBankAccountsRepo(mysql *mysql.MySQL) *BankAccountsRepo {
	return &BankAccountsRepo{mysql}
}

// CreateBankAccount - used to create bank account in the database.
func (b *BankAccountsRepo) CreateBankAccount(ctx context.Context, inp *service.BankAccountInput, client string) (*domain.BankAccount, error) {
	secretValueBytes, err := bcrypt.GenerateFromPassword([]byte(inp.SecretValue), 14)
	if err != nil {
		return nil, err
	}

	card, iban := utils.GenerateNumber(int(inp.ITN))
	account := &domain.BankAccount{
		Client:      client,
		SecretValue: string(secretValueBytes),
		ITN:         inp.ITN,
		CardNumber:  card,
		IBAN:        iban,
		Balance:     0,
	}

	err = b.DB.
		WithContext(ctx).
		Create(account).
		Error
	if err != nil {
		return nil, err
	}

	return account, nil
}

// TopUpBankAccount - used to top up bank account in the database.
func (b *BankAccountsRepo) TopUpBankAccount(ctx context.Context, inp *service.TopUpBankAccountInput, balance float64) error {
	err := b.DB.
		Model(domain.BankAccount{}).
		Where("card_number = ?", inp.CardNumber).
		Update("balance", balance).
		Error
	if err != nil {
		return err
	}

	return err
}

// Check the credit card in the database.
func (b *BankAccountsRepo) CheckCreditCard(ctx context.Context, cardNumber int64) (*domain.BankAccount, error) {
	var card domain.BankAccount
	err := b.DB.
		WithContext(ctx).
		Where("card_number = ?", cardNumber).
		First(&card).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "Card number not found"}
		}
		return nil, err
	}

	return &card, nil
}

// Get credit card info by IBAN in the database.
func (b *BankAccountsRepo) GetInfoByIBAN(ctx context.Context, IBAN string) (*domain.BankAccount, error) {
	var card domain.BankAccount
	err := b.DB.
		WithContext(ctx).
		Where("iban = ?", IBAN).
		First(&card).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "IBAN not found"}
		}
		return nil, err
	}

	return &card, nil
}

// Change the credit card in the database.
func (b *BankAccountsRepo) ChangeCreditCardStatus(ctx context.Context, cardNumber int64, status string) (string, error) {
	err := b.DB.
		Model(domain.BankAccount{}).
		Where("card_number = ?", cardNumber).
		Update("status", status).
		Error
	if err != nil {
		return "", err
	}
	updatedStatus := fmt.Sprintf("Status changed to %s", status)

	return updatedStatus, err
}
