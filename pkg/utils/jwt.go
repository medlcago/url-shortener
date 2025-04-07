package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/google/uuid"
	"net/http"
	"time"
	"url-shortener/internal/models"
)

const (
	defaultDuration = time.Minute * 60
)

type Claims struct {
	Email string    `json:"email"`
	ID    uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *models.User, secret string, duration ...time.Duration) (string, error) {
	tokenDuration := defaultDuration
	if len(duration) > 0 {
		tokenDuration = duration[0]
	}

	claims := &Claims{
		Email: user.Email,
		ID:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractClaimsFromRequest(r *http.Request, secret string) (*Claims, error) {
	tokenString, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		return nil, err
	}

	claims, err := ExtractClaimsFromToken(tokenString, secret)
	return claims, err
}

func ExtractClaimsFromToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
