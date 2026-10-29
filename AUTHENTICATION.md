# Authentication System - Implementation Complete

## ✅ What Was Added

Real authentication has been fully implemented with database integration.

---

## 🔐 Demo User Accounts

All demo users have the same password for testing:

### **Password:** `password123`

### Available Accounts:

| Email | Role | Password |
|-------|------|----------|
| `user@example.com` | User (Student) | `password123` |
| `student@campus.edu` | User (Student) | `password123` |
| `responder@campus.edu` | Responder (Security) | `password123` |
| `warden@hostel.edu` | Responder (Warden) | `password123` |
| `admin@campus.edu` | Admin | `password123` |

---

## 📝 Architecture

### New Files Created:

```
backend/
├── internal/
│   ├── models/
│   │   └── user.go              ← User data models + requests/responses
│   ├── user/
│   │   ├── repository.go        ← Database operations
│   │   └── service.go           ← Business logic (login/signup)
│   └── api/
│       └── router.go            ← Updated with real handlers
├── migrations/
│   ├── 001_init.sql             ← Database schema
│   └── 002_seed_demo_users.sql  ← Demo user data
└── go.mod                        ← Updated with google/uuid dependency
```

### Flow Diagram

```
Frontend (Login Form)
        ↓
  POST /api/auth/login
        ↓
  Backend Handler (router.go)
        ↓
  User Service (login/signup)
        ↓
  User Repository (database queries)
        ↓
  Database (users table)
        ↓
  Auth Service (JWT token generation)
        ↓
  Response: {token, role, name}
```

---

## 🔧 How It Works

### Login Flow

1. User submits email and password
2. Backend queries database for user
3. Password verified (SHA256 hash + JWT secret)
4. JWT token generated with user ID
5. Token, role, and name returned to frontend
6. Frontend stores token in localStorage
7. All subsequent API calls include token in Authorization header

### Signup Flow

1. User provides email, password, name, and role
2. Backend checks if email already exists
3. Password hashed
4. New user created in database
5. JWT token generated
6. Same response as login

---

## 📊 Database Schema

### Users Table
```sql
users {
  id (UUID)           -- Primary key
  email (VARCHAR)     -- Unique, indexed
  password_hash       -- SHA256 hash
  role                -- 'user' | 'responder' | 'admin'
  name                -- Display name
  created_at          -- Timestamp
  updated_at          -- Timestamp
}
```

### Incidents Table
```sql
incidents {
  id (UUID)
  user_id             -- Foreign key to users
  latitude
  longitude
  status              -- 'raised' | 'matched' | 'escalated' | 'resolved'
  message
  created_at
  updated_at
}
```

### Incident Events Table (Audit Trail)
```sql
incident_events {
  id (UUID)
  incident_id         -- Foreign key to incidents
  event_type          -- Event description
  actor_id            -- Who performed action
  detail              -- Additional info
  created_at
}
```

### Responders Table
```sql
responders {
  id (UUID)
  user_id             -- Foreign key to users (responder)
  latitude            -- Current location
  longitude
  is_active           -- Availability flag
  last_heartbeat      -- Last seen timestamp
  created_at
}
```

---

## 🚀 Deployment Instructions

### Step 1: Create Database

```bash
# On Railway PostgreSQL:
psql $POSTGRES_DSN < backend/migrations/001_init.sql
psql $POSTGRES_DSN < backend/migrations/002_seed_demo_users.sql
```

Or use Railway dashboard:
1. Go to your Postgres service
2. Click "Connect"
3. Run the SQL files in the Query editor

### Step 2: Set Environment Variables in Vercel

```
POSTGRES_DSN=postgresql://user:password@host/dbname
REDIS_ADDR=redis-host:6379
REDIS_PASSWORD=***
JWT_SECRET=your-production-secret-key
SOS_ACK_SLA=60s
DISPATCH_RADIUS_KM=2.0
HTTP_PORT=8080
```

### Step 3: Deploy

Push to GitHub:
```bash
git add -A
git commit -m "feat: add real authentication system"
git push origin main
```

Vercel auto-deploys when it detects changes.

---

## 🧪 Testing

### Test Login with curl

```bash
# Successful login
curl -X POST https://your-domain/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "role": "user"
  }'

# Response:
{
  "token": "eyJhbGc...",
  "role": "user",
  "name": "Demo User"
}
```

### Test Signup

```bash
curl -X POST https://your-domain/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123",
    "name": "New User",
    "role": "user"
  }'
```

---

## ⚠️ Security Notes

### Current Implementation (MVP)
- ✅ Passwords hashed with SHA256 + JWT secret
- ✅ JWT tokens for session management
- ✅ Password validation on login
- ✅ Email uniqueness constraint

### Production Improvements Needed
- ❌ Use bcrypt/argon2 instead of SHA256
- ❌ Add rate limiting on login attempts
- ❌ Add HTTPS/TLS enforcement
- ❌ Add password reset functionality
- ❌ Add email verification
- ❌ Add refresh token rotation
- ❌ Add session expiration
- ❌ Add audit logging

---

## 🎯 Next Steps

1. **Database Setup**: Run migrations on your PostgreSQL database
2. **Deploy**: Push to GitHub and redeploy on Vercel
3. **Test**: Use demo credentials to login
4. **Implement**: Add actual SOS dispatch logic, responder matching, etc.
5. **Security**: Add bcrypt, rate limiting, email verification before production

---

## 🆘 Troubleshooting

### "user not found" error
- Ensure demo users were seeded with `002_seed_demo_users.sql`
- Check PostgreSQL is running and connected

### "invalid email or password"
- Verify email and password match exactly
- Note: Password is `password123` for all demo accounts

### Database connection error
- Verify `POSTGRES_DSN` environment variable is set
- Check PostgreSQL service is running on Railway
- Verify firewall/network access

---

## 📚 Files Modified

- `backend/go.mod` - Added google/uuid dependency
- `backend/go.sum` - Updated dependencies
- `backend/internal/api/router.go` - Real login/signup handlers
- `backend/internal/auth/auth.go` - Password hashing + JWT
- New: `backend/internal/models/user.go`
- New: `backend/internal/user/repository.go`
- New: `backend/internal/user/service.go`
- New: `backend/migrations/001_init.sql`
- New: `backend/migrations/002_seed_demo_users.sql`

---

**Ready to deploy!** 🚀 Your authentication system is fully integrated with the database.
