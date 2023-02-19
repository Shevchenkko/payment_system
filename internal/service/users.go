// Package service implements application services.
package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	// external
	"github.com/Shevchenkko/payment_system/pkg/access"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// UsersService - represents users service.
type UsersService struct {
	repos Repositories
	apis  APIs
}

// NewUserService - creates instance of new user service.
func NewUserService(repos Repositories, apis APIs) *UsersService {
	return &UsersService{repos, apis}
}

// RegisterUser is used for creating user.
func (us *UsersService) RegisterUser(ctx context.Context, inp *RegisterUserInput) (RegisterUserOutput, error) {
	// create user in db
	user, err := us.repos.Users.CreateUser(ctx, inp)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	var role access.UserRole
	if user.FullName == os.Getenv("ADMIN_USER_FULLNAME") && user.Email == os.Getenv("ADMIN_USER_EMAIL") {
		role = access.UserRoleAdmin
	} else {
		role = access.UserRoleUser
	}

	// sign auth token
	token, err := access.EncodeToken(
		&access.Token{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(), // 14 days
				NotBefore: time.Now().Unix(),
				Issuer:    "pay-system-api",
				IssuedAt:  time.Now().Unix(),
			},
			UserID:   user.ID,
			UserRole: role,
		},
		os.Getenv("HMAC_SECRET"),
	)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		Token:    token,
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

// LoginUser is used to login a user.
func (us *UsersService) LoginUser(ctx context.Context, inp *LoginUserInput) (LoginUserOutput, error) {
	user, err := us.repos.Users.GetUser(ctx, inp.Email)
	if err != nil {
		return LoginUserOutput{}, err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inp.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return LoginUserOutput{}, &Error{Message: "Wrong password"}
		}

		return LoginUserOutput{}, err
	}

	var role access.UserRole
	if user.FullName == os.Getenv("ADMIN_USER_FULLNAME") && user.Email == os.Getenv("ADMIN_USER_EMAIL") {
		role = access.UserRoleAdmin
	} else {
		role = access.UserRoleUser
	}

	// sign auth token
	t := time.Now()
	token, err := access.EncodeToken(
		&access.Token{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: t.Add(time.Hour * 24 * 14).Unix(), // 14 days
				NotBefore: t.Unix(),
				Issuer:    "pay-system-api",
				IssuedAt:  t.Unix(),
			},
			UserID:   user.ID,
			UserRole: role,
		},
		os.Getenv("HMAC_SECRET"),
	)
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		Token:    token,
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

// VerifyAccessToken is used to verify user jwt access token.
func (us *UsersService) VerifyAccessToken(ctx context.Context, token string) (bool, int, string) {
	tokenData, err := access.DecodeToken(token, os.Getenv("HMAC_SECRET"))
	if err != nil {
		return false, 0, ""
	}

	return true, tokenData.UserID, string(tokenData.UserRole)
}

// SendEmail is used for sending email.
func (us *UsersService) SendEmail(ctx context.Context, inp *SendUserEmailInput) error {
	// get user from db
	user, err := us.repos.Users.GetUser(ctx, inp.Email)
	if err != nil {
		return err
	}

	var role access.UserRole
	if user.FullName == os.Getenv("ADMIN_USER_FULLNAME") && user.Email == os.Getenv("ADMIN_USER_EMAIL") {
		role = access.UserRoleAdmin
	} else {
		role = access.UserRoleUser
	}

	// generate token
	t := time.Now()
	token, err := access.EncodeToken(
		&access.Token{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: t.Add(time.Minute * 15).Unix(), // 15 minute
				NotBefore: t.Unix(),
				Issuer:    inp.Email,
				IssuedAt:  t.Unix(),
			},
			UserID:   user.ID,
			UserRole: role,
		},
		os.Getenv("HMAC_SECRET"),
	)
	if err != nil {
		return err
	}

	// create token in db
	err = us.repos.Users.CreateToken(ctx,
		GenerateTokenInput{
			Email: inp.Email,
			Token: token,
		})
	if err != nil {
		return err
	}

	err = us.apis.Emails.SendEmail(ctx, SendEmailInput{
		To:          inp.Email,
		Subject:     "PaySystem: reset password",
		ContentType: "text/html",
		Body: fmt.Sprintf(`
			<h2>PaySystem: reset password</h1>
	    <p>Hello!</p>
	    <p>Someone (we hope it was you) decided to change the forgotten password for the PaySystem account associated with this email address.<p>
		<p>Please, use this token for to change your account password:</p>
	    <p>token=%s</p>
		`, token),
	})
	if err != nil {
		return err
	}

	return err
}

// ResetPassword is used for reset password.
func (us *UsersService) ResetPassword(ctx context.Context, inp *ResetPasswordInput) error {
	// get token from db
	_, err := us.repos.Users.GetToken(ctx, inp.Token)
	if err != nil {
		return err
	}

	// verify user token
	_, err = access.DecodeToken(inp.Token, os.Getenv("HMAC_SECRET"))
	if err != nil {
		return err
	}

	// change password in db
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(inp.Password), 14)
	if err != nil {
		return err
	}

	err = us.repos.Users.ResetPassword(ctx,
		&ResetPasswordInput{
			Token:    inp.Token,
			Password: string(passwordBytes),
		})
	if err != nil {
		return err
	}

	// delete token in db
	err = us.repos.Users.DeleteToken(ctx, inp.Token)
	if err != nil {
		return err
	}

	return nil
}
