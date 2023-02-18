// Package domain implements application domain.
package domain

import (
	"github.com/Shevchenkko/payment_system/pkg/mysql"
)

// User represents the user model stored in the database.
type User struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	FullName string `json:"fullName,omitempty" binding:"required"`
	Email    string `json:"email,omitempty" gorm:"column:email;not null;unique;index" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`

	mysql.Model
}

// UserToken represents the token model stored in the database.
type UserToken struct {
	ID    int    `json:"id,omitempty" gorm:"primaryKey"`
	Email string `json:"email,omitempty" gorm:"index" binding:"required"`
	Token string `json:"token,omitempty" gorm:"index" binding:"required"`

	mysql.Model
}
