package user

import "github.com/rishabhjha/rakshagrid/internal/models"

// UserService interface for both database and demo implementations
type UserService interface {
	Login(email, password string) (token string, role string, userName string, err error)
	Signup(email, password, name, role string) (token string, resultRole string, resultName string, err error)
	GetUserByID(id string) (*models.User, error)
}

// Ensure both services implement the interface
var _ UserService = (*Service)(nil)
var _ UserService = (*DemoService)(nil)
