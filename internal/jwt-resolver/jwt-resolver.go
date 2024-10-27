package jwtresolver

import (
	"errors"
	"time"

	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Settings used to initialize JWTResolver.
type Settings struct {
	// Secret used to sign and verify JWT tokens
	// Must be at least 32 bytes
	Secret string

	// Issuer is token issuer
	Issuer string

	// TokenTTL is token lifetime
	TokenTTL time.Duration
}

// JWTResolver used to create and decrypt JWT tokens.
// It uses HS256 signing method.
// Must be initialized with New method.
type JWTResolver struct {
	secret   string
	issuer   string
	tokenTTL time.Duration
}

// New returns new JWTResolver instance.
func New(settings Settings) *JWTResolver {
	return &JWTResolver{
		secret:   settings.Secret,
		issuer:   settings.Issuer,
		tokenTTL: settings.TokenTTL,
	}
}

// DecryptToken decodes JWT token and returns it.
// Returns error if token decoding failed.
func (res *JWTResolver) DecryptToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error unexpected signing method")
		}
		return []byte(res.secret), nil
	})

	// Check token validity
	if !token.Valid {
		return claims, errors.New("error invalid token")
	}

	return claims, err
}

// CreateToken creates JWT token and returns it.
// Returns error if token creation failed.
func (res *JWTResolver) CreateToken(subject string, userID uuid.UUID) (string, error) {
	// Create token
	claims := models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    res.issuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(res.tokenTTL)),
		},
		UserID: userID,
	}

	// Sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(res.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
