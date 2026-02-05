package jwt

import (
	"errors"
	"fmt"
	"meye-core/internal/domain/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var _ user.JWTService = (*Service)(nil)

type Service struct {
	secret         string
	issuer         string
	expirationTime time.Duration
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func New(secret, issuer string, expirationTime time.Duration) *Service {
	return &Service{
		secret:         secret,
		issuer:         issuer,
		expirationTime: expirationTime,
	}
}

func (s *Service) GenerateSignedToken(user *user.User) (string, error) {
	claims := Claims{
		UserID: user.ID(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   user.ID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses a token string and validates it, returning the identification claims.
func (s *Service) ValidateToken(tokenString string) (string, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token is not valid")
	}

	return claims.UserID, nil
}
