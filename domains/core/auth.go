package core

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// SigningSecret global signing secret variable
var SigningSecret string

// TokenData returns the base information storeed within a JWT token
type TokenData struct {
	UserID, Email, UserType string
}

// AuthClaims extra claims struct for using standard claims
type AuthClaims struct {
	*TokenData
	jwt.RegisteredClaims
}

// Generate generates an auth secret
func (t *TokenData) Generate() (string, error) {
	claims := &AuthClaims{
		t,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // make it expire in a month
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "subs-manager",
			Subject:   "auth-token",
			ID:        t.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SigningSecret))
}

// Decode decodes the jwt value and returns token data
func Decode(token string) (*TokenData, error) {
	tk, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(SigningSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tk.Claims.(*AuthClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims.TokenData, nil
}

// GeneratePassword generates password from the value user provides
func GeneratePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash comparse the provided hash and the user provided password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
