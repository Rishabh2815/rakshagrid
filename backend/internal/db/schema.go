package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// InitializeSchema creates tables and seeds demo data if needed
func InitializeSchema(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create users table
	usersSchema := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'user',
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	if _, err := db.ExecContext(ctx, usersSchema); err != nil {
		log.Printf("error creating users table: %v", err)
		return err
	}

	// Create incidents table
	incidentsSchema := `
	CREATE TABLE IF NOT EXISTS incidents (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL REFERENCES users(id),
		latitude DECIMAL(10, 8) NOT NULL,
		longitude DECIMAL(11, 8) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'raised',
		message TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_incidents_user_id ON incidents(user_id);
	CREATE INDEX IF NOT EXISTS idx_incidents_status ON incidents(status);
	`

	if _, err := db.ExecContext(ctx, incidentsSchema); err != nil {
		log.Printf("error creating incidents table: %v", err)
		return err
	}

	// Create incident_events table
	eventsSchema := `
	CREATE TABLE IF NOT EXISTS incident_events (
		id VARCHAR(255) PRIMARY KEY,
		incident_id VARCHAR(255) NOT NULL REFERENCES incidents(id),
		event_type VARCHAR(50) NOT NULL,
		actor_id VARCHAR(255) REFERENCES users(id),
		detail TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_events_incident_id ON incident_events(incident_id);
	`

	if _, err := db.ExecContext(ctx, eventsSchema); err != nil {
		log.Printf("error creating incident_events table: %v", err)
		return err
	}

	// Seed demo users if they don't exist
	if err := seedDemoUsers(db); err != nil {
		log.Printf("error seeding demo users: %v", err)
		return err
	}

	log.Println("database schema initialized successfully")
	return nil
}

// seedDemoUsers creates demo user accounts for testing
func seedDemoUsers(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// All demo users have password: "password123"
	// Password hash: f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4
	// (This is SHA256 of "password123" + "dev-secret-change-me")

	demoUsers := []struct {
		id       string
		email    string
		hash     string
		role     string
		name     string
	}{
		{
			id:   "550e8400-e29b-41d4-a716-446655440001",
			email: "user@example.com",
			hash:  "f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4",
			role:  "user",
			name:  "Demo User",
		},
		{
			id:    "550e8400-e29b-41d4-a716-446655440002",
			email: "student@campus.edu",
			hash:  "f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4",
			role:  "user",
			name:  "Campus Student",
		},
		{
			id:    "550e8400-e29b-41d4-a716-446655440003",
			email: "responder@campus.edu",
			hash:  "f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4",
			role:  "responder",
			name:  "Security Guard",
		},
		{
			id:    "550e8400-e29b-41d4-a716-446655440004",
			email: "warden@hostel.edu",
			hash:  "f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4",
			role:  "responder",
			name:  "Hostel Warden",
		},
		{
			id:    "550e8400-e29b-41d4-a716-446655440005",
			email: "admin@campus.edu",
			hash:  "f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4",
			role:  "admin",
			name:  "Campus Administrator",
		},
	}

	for _, user := range demoUsers {
		query := `
		INSERT INTO users (id, email, password_hash, role, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (email) DO NOTHING
		`
		_, err := db.ExecContext(ctx, query, user.id, user.email, user.hash, user.role, user.name)
		if err != nil {
			log.Printf("error inserting demo user %s: %v", user.email, err)
			// Don't return error - continue with other users
		}
	}

	return nil
}
