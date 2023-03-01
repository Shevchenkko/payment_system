package domain

import "github.com/Shevchenkko/payment_system/pkg/mysql"

// MessageLog represents the message log model stored in the database.
type MessageLog struct {
	ID      int    `json:"id,omitempty" gorm:"primaryKey"`
	Client  string `json:"client,omitempty" gorm:"column:client"`
	Message string `json:"message,omitempty" gorm:"column:message"`

	mysql.Model
}
