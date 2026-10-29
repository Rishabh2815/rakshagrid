package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	jwtSecret string
	jwtExpiry time.Duration
}

func NewService(secret string, expiry time.Duration) *Service {
	return &Service{
		jwtSecret: secret,
		jwtExpiry: expiry,
	}
}

// HashPassword creates a hash of the password (MVP - use bcrypt in production)
func (s *Service) HashPassword(password string) string {
	// NOTE: This is MVP only. Use bcrypt/argon2 in production.
	h := sha256.New()
	h.Write([]byte(password + s.jwtSecret))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// VerifyPassword checks if a password matches its hash
func (s *Service) VerifyPassword(password, hash string) bool {
	return s.HashPassword(password) == hash
}

// GenerateToken creates a JWT token for a user
func (s *Service) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// VerifyToken validates a JWT token and returns the user ID
func (s *Service) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
	}
	return "", fmt.Errorf("invalid token")
}
