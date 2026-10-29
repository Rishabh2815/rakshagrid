package user

import (
	"errors"

	"github.com/rishabhjha/rakshagrid/internal/auth"
	"github.com/rishabhjha/rakshagrid/internal/models"
)

type Service struct {
	repo      *Repository
	authSvc   *auth.Service
}

func NewService(repo *Repository, authSvc *auth.Service) *Service {
	return &Service{
		repo:    repo,
		authSvc: authSvc,
	}
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(email, password string) (string, string, string, error) {
	// Get user from database
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", "", errors.New("invalid email or password")
	}

	// Verify password
	if !s.authSvc.VerifyPassword(password, user.PasswordHash) {
		return "", "", "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return "", "", "", err
	}

	return token, user.Role, user.Name, nil
}

// Signup creates a new user account
func (s *Service) Signup(email, password, name, role string) (string, string, string, error) {
	// Check if email already exists
	_, err := s.repo.GetUserByEmail(email)
	if err == nil {
		return "", "", "", errors.New("email already registered")
	}

	// Hash password
	passwordHash := s.authSvc.HashPassword(password)

	// Create new user
	user := models.NewUser(email, passwordHash, role, name)

	// Save to database
	if err := s.repo.CreateUser(user); err != nil {
		return "", "", "", err
	}

	// Generate JWT token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return "", "", "", err
	}

	return token, user.Role, user.Name, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetUserByID(id)
}
