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

// Search users - used to search user from the database.
func (r *UsersRepo) SearchUsers(ctx context.Context, filter *domain.Filter) (*service.SearchUsers, error) {
	if filter == nil {
		filter = new(domain.Filter)
		filter.Validate()
	}

	q := r.DB.
		Table("users").
		Offset((filter.Page - 1) * filter.List).
		Limit(filter.List).
		Order(filter.OrderString())

	var userOutput []service.User
	var response *service.SearchUsers
	if err := q.Find(&userOutput).Error; err != nil {
		return nil, &service.Error{Message: "Users not found"}
	}

	var count int64
	q = r.DB.
		Table("users")
	if err := q.Count(&count).Error; err != nil {
		return nil, &service.Error{Message: "Users not found"}
	}

	response = &service.SearchUsers{
		Data: userOutput,
		Pagination: &domain.Pagination{
			Order: filter.OrderString(),
			Page:  filter.Page,
			List:  filter.List,
			Total: &count,
		},
	}

	return response, nil
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

// GetUserByID is used to get user by id from the database.
func (r *UsersRepo) GetUserByID(ctx context.Context, userId int) (*domain.User, error) {
	var user domain.User
	err := r.DB.
		WithContext(ctx).
		Where("id = ?", userId).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "User not found"}
		}
		return nil, err
	}

	return &user, err
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
