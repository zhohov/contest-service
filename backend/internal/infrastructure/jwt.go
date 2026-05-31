package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zhohov/contest-service/backend/internal/domain"
)

var ErrInvalidToken = errors.New("Invalid token")

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JwtService struct {
	secret []byte
	ttl    time.Duration
}

func NewJwtService(secret []byte, ttl time.Duration) *JwtService {
	return &JwtService{
		secret: secret,
		ttl:    ttl,
	}
}

func (s *JwtService) GenerateToken(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(s.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JwtService) ParseToken(token string) (*domain.AuthUser, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &domain.AuthUser{
		ID:   int(userID),
		Role: role,
	}, nil
}
