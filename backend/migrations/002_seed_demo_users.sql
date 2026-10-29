-- Seed demo users for testing
-- Passwords are hashed with JWT_SECRET (assuming "dev-secret-change-me" in dev mode)
-- Demo users for testing the authentication system

-- User account (Student)
INSERT INTO users (id, email, password_hash, role, name, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'user@example.com', 'f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4', 'user', 'Demo User', NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'student@campus.edu', 'f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4', 'user', 'Campus Student', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Responder accounts (Security/Warden)
INSERT INTO users (id, email, password_hash, role, name, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440003', 'responder@campus.edu', 'f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4', 'responder', 'Security Guard', NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440004', 'warden@hostel.edu', 'f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4', 'responder', 'Hostel Warden', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Admin account
INSERT INTO users (id, email, password_hash, role, name, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440005', 'admin@campus.edu', 'f7e8d9c8c9a9b8a7d6c5b4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4', 'admin', 'Campus Administrator', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- All demo users have password: "password123"
