// Package access implements server access.
package access

import (
	"fmt"

	// third party
	"github.com/dgrijalva/jwt-go"
)

// Token represents authentication token struct.
// Used for protecting private endpoints.
type Token struct {
	jwt.StandardClaims

	UserID   int      `json:"userId"`
	UserRole UserRole `json:"userRole"`
}

// Role represents a user role.
type UserRole string

// RoleUser represents a plain user role.
var UserRoleUser UserRole = "user"

// RoleAdmin represents a plain user role.
var UserRoleAdmin UserRole = "admin"

// EncodeToken is used to encode Token to string.
func EncodeToken(token *Token, hmacSecret string) (string, error) {
	// create jwt token object
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)

	// encode jwt object to string
	tokenString, err := jwtToken.SignedString([]byte(hmacSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// DecodeToken is used to verify and decode Token from string.
func DecodeToken(tokenString string, hmacSecret string) (*Token, error) {
	// parse token
	jwtToken, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		// validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// return validation secret
		return []byte(hmacSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("malformed token")
	}

	// check if token is valid
	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// parse claims
	token, ok := jwtToken.Claims.(*Token)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return token, nil
}
