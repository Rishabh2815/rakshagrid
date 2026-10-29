package user

import (
	"errors"
	"sync"

	"github.com/rishabhjha/rakshagrid/internal/auth"
	"github.com/rishabhjha/rakshagrid/internal/models"
)

// DemoService provides in-memory user authentication for demo/testing
type DemoService struct {
	users map[string]*models.User
	mu    sync.RWMutex
	authSvc *auth.Service
}

func NewDemoService(authSvc *auth.Service) *DemoService {
	svc := &DemoService{
		users:   make(map[string]*models.User),
		authSvc: authSvc,
	}
	svc.seedDemoUsers()
	return svc
}

func (s *DemoService) seedDemoUsers() {
	// Demo users - password: "password123"
	passwordHash := s.authSvc.HashPassword("password123")

	demoUsers := []*models.User{
		{
			ID:           "550e8400-e29b-41d4-a716-446655440001",
			Email:        "user@example.com",
			PasswordHash: passwordHash,
			Role:         "user",
			Name:         "Demo User",
		},
		{
			ID:           "550e8400-e29b-41d4-a716-446655440002",
			Email:        "student@campus.edu",
			PasswordHash: passwordHash,
			Role:         "user",
			Name:         "Campus Student",
		},
		{
			ID:           "550e8400-e29b-41d4-a716-446655440003",
			Email:        "responder@campus.edu",
			PasswordHash: passwordHash,
			Role:         "responder",
			Name:         "Security Guard",
		},
		{
			ID:           "550e8400-e29b-41d4-a716-446655440004",
			Email:        "warden@hostel.edu",
			PasswordHash: passwordHash,
			Role:         "responder",
			Name:         "Hostel Warden",
		},
		{
			ID:           "550e8400-e29b-41d4-a716-446655440005",
			Email:        "admin@campus.edu",
			PasswordHash: passwordHash,
			Role:         "admin",
			Name:         "Campus Administrator",
		},
	}

	for _, user := range demoUsers {
		s.users[user.Email] = user
	}
}

// Login authenticates a user against demo data
func (s *DemoService) Login(email, password string) (string, string, string, error) {
	s.mu.RLock()
	user, exists := s.users[email]
	s.mu.RUnlock()

	if !exists {
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

// Signup creates a new user in demo mode
func (s *DemoService) Signup(email, password, name, role string) (string, string, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if email already exists
	if _, exists := s.users[email]; exists {
		return "", "", "", errors.New("email already registered")
	}

	// Hash password
	passwordHash := s.authSvc.HashPassword(password)

	// Create new user
	user := models.NewUser(email, passwordHash, role, name)

	// Store in memory
	s.users[email] = user

	// Generate JWT token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return "", "", "", err
	}

	return token, user.Role, user.Name, nil
}

// GetUserByID retrieves a user by ID
func (s *DemoService) GetUserByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}
