package domain

import "github.com/Shevchenkko/payment_system/pkg/mysql"

// BankAccount represents the bank account model stored in the database.
type BankAccount struct {
	ID           int     `json:"id,omitempty" gorm:"primaryKey"`
	Client       string  `json:"client,omitempty" gorm:"column:client"`
	ClientStatus string  `json:"clientStatus,omitempty" gorm:"column:client_status;type:enum('ACTIVE','LOCK');default:'ACTIVE'"`
	SecretValue  string  `json:"secretValue" gorm:"column:secret_value"`
	ITN          int64   `json:"itn,omitempty" gorm:"column:itn;not null;index"`
	CardNumber   int64   `json:"cardNumber,omitempty" gorm:"column:card_number;not null;unique;index"`
	IBAN         string  `json:"iban,omitempty" gorm:"column:iban;not null;unique;index"`
	Balance      float64 `json:"balance,omitempty" gorm:"column:balance"`
	Status       string  `json:"status,omitempty" gorm:"column:status;type:enum('ACTIVE','LOCK');default:'ACTIVE'"`

	mysql.Model
}
