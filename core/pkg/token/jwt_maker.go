package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secret string
}

var (
	managerInstance *JWTManager
)

type TokensAndClaims struct {
	AccessToken        string
	AccessTokenClaims  *UserClaims
	RefreshToken       string
	RefreshTokenClaims *UserClaims
}

func NewJWTManager(secret string) *JWTManager {
	managerInstance = &JWTManager{secret}
	return managerInstance
}

func GetManager() *JWTManager {
	if managerInstance == nil {
		panic("Token manager instace doesn't exist")
	}

	return managerInstance
}

var (
	ErrInvalidSigningMethod = errors.New("invalid token signing method")
	ErrInvalidClaims        = errors.New("invalid token claims")
	ErrTokenExpired         = errors.New("token is expired")
)

func (m *JWTManager) CreateToken(id uuid.UUID, email string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (m *JWTManager) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(m.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, fmt.Errorf("error in parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}
