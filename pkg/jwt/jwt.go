package jwt

import (
	"errors"
	jwtToken "github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"net/http"
	"time"
)

const (
	accessDefaultDuration  = time.Minute * 60
	refreshDefaultDuration = time.Hour * 24 * 30
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

type Claims struct {
	ID   string    `json:"id"`
	Type TokenType `json:"token_type"`
	jwtToken.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWT(userID, secret string, tokenType TokenType, duration time.Duration) (string, error) {
	claims := &Claims{
		ID:   userID,
		Type: tokenType,
		RegisteredClaims: jwtToken.RegisteredClaims{
			ExpiresAt: jwtToken.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwtToken.NewWithClaims(jwtToken.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewAccessToken(userID, secret string, duration ...time.Duration) (string, error) {
	tokenDuration := accessDefaultDuration
	if len(duration) > 0 {
		tokenDuration = duration[0]
	}
	return NewJWT(userID, secret, Access, tokenDuration)
}

func NewRefreshToken(userID, secret string, duration ...time.Duration) (string, error) {
	tokenDuration := refreshDefaultDuration
	if len(duration) > 0 {
		tokenDuration = duration[0]
	}
	return NewJWT(userID, secret, Refresh, tokenDuration)
}

func NewTokenPair(userID, secret string) (*TokenPair, error) {
	accessToken, err := NewAccessToken(userID, secret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := NewRefreshToken(userID, secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ExtractClaimsFromRequest(r *http.Request, secret string) (*Claims, error) {
	tokenString, err := ExtractTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	claims, err := ExtractClaimsFromToken(tokenString, secret)
	return claims, err
}

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	tokenString, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractClaimsFromToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwtToken.ParseWithClaims(tokenString, claims, func(token *jwtToken.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwtToken.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
