package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/e1m0re/passman/internal/model"
)

type JWTManager interface {
	// Generate creates new JWT-token.
	Generate(user *model.User) (string, error)
	// Verify validates access token.
	Verify(accessToken string) (*UserClaims, error)
}

type jwtManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

// Generate creates new JWT-token.
func (m *jwtManager) Generate(user *model.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
		},
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString([]byte(m.secretKey))
}

// Verify validates access token.
func (m *jwtManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return []byte(m.secretKey), nil
		})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

var _ JWTManager = (*jwtManager)(nil)

// NewJWTManager initiates new instance of jwtManager.
func NewJWTManager(secretKey string, tokenDuration time.Duration) JWTManager {
	return &jwtManager{secretKey, tokenDuration}
}
