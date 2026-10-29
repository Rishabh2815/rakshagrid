# Vercel Services Configuration - Review & Confirmation

Your RakshaGrid repository has been configured for **Vercel Services** deployment. Please review and confirm the configuration details below.

## Final Configuration Summary

### Service Names ✅

```json
{
  "services": {
    "frontend": { "root": "frontend", "framework": "vite" },
    "backend": { "root": "backend", "runtime": "go" }
  }
}
```

- **Frontend Service:** `frontend` (React + Vite)
- **Backend Service:** `backend` (Go)

### Public Routes & Routing ✅

| Route | Service | Public? | Purpose |
|-------|---------|---------|---------|
| `/` (catch-all) | frontend | ✅ Yes | React single-page app |
| `/api/*` | backend | ✅ Yes | API endpoints & WebSocket |
| `/health` | backend | ✅ Yes | Health check (at `/api/health` after rewrite) |

**Vercel Rewrites:**
```json
{
  "source": "/api/(.*)",
  "destination": { "service": "backend" }
},
{
  "source": "/(.*)",
  "destination": { "service": "frontend" }
}
```

**Result:** Single unified domain with:
- `https://rakshagrid.vercel.app/` → Frontend
- `https://rakshagrid.vercel.app/api/*` → Backend
- No CORS issues (same-origin requests)

### Service Bindings ✅

**No bindings needed** because:

1. **Frontend does NOT call backend directly**
   - Frontend lives at `/`, makes requests to `/api/`
   - Vercel rewrites `/api/` to backend automatically
   - No service-to-service communication

2. **Backend does NOT call frontend**
   - Backend is API-only, no frontend needs

3. **Communication Flow:**
   ```
   Browser → fetch("/api/sos/raise")
   ↓
   Vercel rewrites /api/* to backend service
   ↓
   Backend handler processes request
   ↓
   Response sent to browser
   ```

**Therefore:** Zero bindings needed. Simple architecture.

---

## Code Changes Made

### 1. Frontend API Client (`frontend/src/api.js`)

**Before:**
```javascript
const API_BASE = process.env.REACT_APP_API_URL || "/v1";
// Called /v1/auth/login, /v1/sos/raise, etc.
```

**After:**
```javascript
const API_BASE = "/api";
// Calls /api/auth/login, /api/sos/raise, etc.
// Vercel rewrites /api/* to backend
```

✅ **Change:** All API calls now use `/api/` prefix (no env var needed)

### 2. Backend Router (`backend/internal/api/router.go`)

**Before:**
```go
r.Route("/v1", func(r chi.Router) {
  r.Post("/sos/raise", handleRaiseSOS(...))
})
// Served at /v1/sos/raise
```

**After:**
```go
r.Route("/auth", func(r chi.Router) {
  r.Post("/login", handleLogin(...))
})
r.Route("/sos", func(r chi.Router) {
  r.Post("/raise", handleRaiseSOS(...))
})
// Served at /auth/login, /sos/raise
// But Vercel serves these at /api/auth/login, /api/sos/raise
```

✅ **Change:** Removed `/v1` prefix (backend routes are now `/auth`, `/sos`, etc.)

### 3. Frontend Vite Config (`frontend/vite.config.js`)

**Before:**
```javascript
proxy: {
  "/v1": {
    target: process.env.VITE_BACKEND_URL || "http://localhost:8080"
  }
}
```

**After:**
```javascript
proxy: {
  "/api": {
    target: "http://localhost:8080"
  }
}
```

✅ **Change:** Dev proxy routes `/api` to backend (for local testing)

### 4. Backend Config (`backend/internal/config/config.go`)

**Before:**
```go
AllowedOrigins: []string{
  "http://localhost:3000",
  "http://localhost:5173",
  "https://localhost",
  os.Getenv("FRONTEND_URL"),
}
```

**After:**
```go
AllowedOrigins: []string{
  "http://localhost:3000",
  "http://localhost:5173",
  "https://localhost",
  // In Vercel services: same-origin (no CORS policy needed)
}
```

✅ **Change:** Removed FRONTEND_URL env var (not needed for same-origin requests)

### 5. New Configuration File (`vercel.json`)

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "services": {
    "frontend": {
      "root": "frontend",
      "framework": "vite",
      "buildCommand": "npm run build",
      "outputDirectory": "dist"
    },
    "backend": {
      "root": "backend",
      "runtime": "go"
    }
  },
  "rewrites": [
    {
      "source": "/api/(.*)",
      "destination": { "service": "backend" }
    },
    {
      "source": "/(.*)",
      "destination": { "service": "frontend" }
    }
  ]
}
```

✅ **New:** Vercel auto-detects this file and configures services accordingly

---

## API Endpoints

All backend routes are now prefixed with `/api/`:

| Endpoint | Method | Service | Purpose |
|----------|--------|---------|---------|
| `GET /health` | GET | backend | Backend health check |
| `POST /api/auth/login` | POST | backend | User login |
| `POST /api/sos/raise` | POST | backend | Raise emergency SOS |
| `GET /api/responders/nearby` | GET | backend | Find nearby responders |
| `GET /api/incidents` | GET | backend | Fetch incidents |
| `WS /api/stream` | WebSocket | backend | Real-time updates |

---

## Deployment Flow

1. **Push to GitHub** ✅ (already done)
2. **Import on Vercel:**
   - Visit https://vercel.com/dashboard
   - "Add New" → "Project"
   - Select `Rishabh2815/rakshagrid`
   - Vercel auto-detects `vercel.json`
   - Configure env vars (if using external DB)
   - Deploy
3. **Result:** Single domain serving both frontend and backend

---

## Local Testing (Before Deploying)

```bash
# Terminal 1: Start backend
cd backend
HTTP_PORT=8080 JWT_SECRET="dev" go run .

# Terminal 2: Start frontend (Vite dev server)
cd frontend
npm run dev
# Opens http://localhost:5173
# Vite proxies /api to http://localhost:8080
```

Test at `http://localhost:5173`:
- Frontend loads ✅
- Click SOS button → POST /api/sos/raise → backend responds ✅

---

## Comparison with Original Setup

| Aspect | Original (Railway) | New (Vercel Services) |
|--------|-------------------|----------------------|
| Frontend URL | `https://rakshagrid.vercel.app` | `https://rakshagrid.vercel.app` |
| Backend URL | `https://rakshagrid-api.railway.app` | `https://rakshagrid.vercel.app/api` |
| Domain Count | 2 domains | 1 domain ✅ |
| CORS Config | Required | Not needed ✅ |
| Databases | Railway Postgres/Redis | Use external services + env vars |
| Cost | $5–10/month (Railway) | Free tier available ✅ |

---

## ⚠️ Please Confirm

Before proceeding, please verify and confirm:

### ✅ Service Names
- **Frontend:** `frontend`
- **Backend:** `backend`
- Correct? **YES** / **NO**

### ✅ Public Routes
- **Frontend on `/`:** Yes (catch-all)
- **Backend on `/api/*`:** Yes
- Correct? **YES** / **NO**

### ✅ Bindings
- **Frontend → Backend:** No direct calls (Vercel rewrites handle routing)
- **Backend → Frontend:** Not needed
- Correct? **YES** / **NO**

### ✅ API Path Change
- **Old path:** `/v1/sos/raise`
- **New path:** `/api/sos/raise`
- Acceptable? **YES** / **NO**

### ✅ Code Changes
All of the following have been updated to use `/api/`:
- `frontend/src/api.js` ✅
- `backend/internal/api/router.go` ✅
- `frontend/vite.config.js` ✅
- All confirmed? **YES** / **NO**

---

## Next Steps

Once you confirm the above, your setup is ready:

1. ✅ All code committed to GitHub
2. Deploy on Vercel using your dashboard
3. Set environment variables (POSTGRES_DSN, REDIS_ADDR, etc.) if using external DB
4. Test the deployment

---

**Please reply with confirmation or any corrections needed.**
