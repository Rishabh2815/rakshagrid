package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rishabhjha/rakshagrid/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a new user into the database
func (r *Repository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.Role, user.Name, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return errors.New("email already registered")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, role, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, role, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
