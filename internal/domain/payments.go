package domain

import (
	"github.com/Shevchenkko/payment_system/pkg/mysql"
)

// Payment represents the payments model stored in the database.
type Payment struct {
	ID                   int64   `json:"id,omitempty" gorm:"primaryKey"`
	PaymentStatus        string  `json:"paymentStatus,omitempty" gorm:"column:payment_status;type:enum('prepared','sent');default:'prepared'"`
	FromClient           string  `json:"fromClient,omitempty" gorm:"column:from_client"`
	FromClientITN        int64   `json:"fromClientItn,omitempty" gorm:"column:from_client_itn;not null;index"`
	FromClientIBAN       string  `json:"fromClientIban,omitempty" gorm:"column:from_client_iban;not null;index"`
	FromClientCardNumber int64   `json:"fromClientCardNumber,omitempty" gorm:"column:from_client_card_number;not null;index"`
	Description          string  `json:"description" gorm:"column:description"`
	ToClientIBAN         string  `json:"toClientIban,omitempty" gorm:"column:to_client_iban;not null;index"`
	ToClient             string  `json:"toClient,omitempty" gorm:"column:to_client"`
	OperationAmount      float64 `json:"operationAmount,omitempty" gorm:"column:operation_amount"`

	mysql.Model
}
