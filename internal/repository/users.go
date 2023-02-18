// Package repository implements application repository.
package repository

import (
	"context"
	"errors"

	// external
	"github.com/Shevchenkko/payment_system/pkg/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	// internal
	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/service"
)

// UsersRepo - represents user repository.
type UsersRepo struct {
	*mysql.MySQL
}

// NewUsersRepo - create new instance of users repo.
func NewUsersRepo(mysql *mysql.MySQL) *UsersRepo {
	return &UsersRepo{mysql}
}

// CreateUser - used to create user in the database.
func (r *UsersRepo) CreateUser(ctx context.Context, inp *service.RegisterUserInput) (*domain.User, error) {
	isExists, _ := r.GetUser(ctx, inp.Email)
	if isExists != nil {
		return nil, &service.Error{Message: "User with this email already exists"}
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(inp.Password), 14)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FullName: inp.FullName,
		Email:    inp.Email,
		Password: string(passwordBytes),
	}

	err = r.DB.
		WithContext(ctx).
		Create(user).
		Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser is used to get a user from the database.
func (r *UsersRepo) GetUser(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "User not registered"}
		}

		return nil, err
	}

	return &user, nil
}

// CreateToken - used to create token in the database.
func (r *UsersRepo) CreateToken(ctx context.Context, inp service.GenerateTokenInput) error {
	err := r.DB.
		WithContext(ctx).
		Create(&domain.UserToken{
			Email: inp.Email,
			Token: inp.Token,
		}).
		Error
	if err != nil {
		return err
	}

	return nil
}

// GetToken is used to get a token from the database.
func (r *UsersRepo) GetToken(ctx context.Context, token string) (*domain.UserToken, error) {
	// var user domain.User
	var user domain.UserToken
	err := r.DB.
		Model(domain.UserToken{}).
		Where("token IN (?)", token).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "User with provided token not found"}
		}

		return nil, err
	}

	return &user, nil
}

// DeleteToken - used to delete token in the database.
func (r *UsersRepo) DeleteToken(ctx context.Context, token string) error {
	err := r.DB.
		Delete(&domain.UserToken{}, "token = ?", token).
		Error
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword is used to reset password from the database.
func (r *UsersRepo) ResetPassword(ctx context.Context, inp *service.ResetPasswordInput) error {
	user, err := r.GetToken(ctx, inp.Token)
	if err != nil {
		return err
	}

	err = r.DB.Model(domain.User{}).
		Where("email IN (?)", user.Email).
		Update("password", inp.Password).
		Error
	if err != nil {
		return err
	}

	return err
}
