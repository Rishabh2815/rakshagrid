# Vercel Services Deployment Configuration

This repository is now configured for deployment as a **Vercel Services** project, combining frontend and backend into a single unified deployment.

## Architecture

```
┌─────────────────────────────────────────────────────┐
│         Vercel Services Project                      │
│      (rakshagrid.vercel.app)                         │
├─────────────────────────────────────────────────────┤
│                                                      │
│  ┌──────────────┐      ┌──────────────┐            │
│  │   Frontend   │      │   Backend    │            │
│  │  (Vite React)│      │  (Go HTTP)   │            │
│  │   on /       │      │   on /api    │            │
│  └──────────────┘      └──────────────┘            │
│                                                      │
│  Requests routed by vercel.json rewrites            │
│  - /                   → Frontend (React)           │
│  - /api/*              → Backend (Go)               │
│  - /api/stream (WS)    → Backend (WebSocket)        │
│                                                      │
└─────────────────────────────────────────────────────┘
```

## Configuration Files

### `vercel.json` (Root)

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

**What this does:**
- Defines two services: `frontend` (Vite/React) and `backend` (Go)
- Routes `/api/*` requests to the Go backend
- Routes all other requests to React frontend (single-page app)
- Both services share the same domain

## API Routes

All API endpoints are served under `/api/`:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Backend health check |
| `/api/auth/login` | POST | User authentication |
| `/api/sos/raise` | POST | Raise emergency SOS |
| `/api/responders/nearby` | GET | Find nearby responders (geospatial) |
| `/api/incidents` | GET | Fetch incidents |
| `/api/stream` | WebSocket | Real-time updates |

## Frontend-Backend Communication

### No Environment Variables Needed

In production (Vercel), API calls automatically use `/api/`:

```javascript
// frontend/src/api.js
const API_BASE = "/api";  // Hardcoded to /api

// All calls go through Vercel rewrites:
fetch(`${API_BASE}/sos/raise`)  // → /api/sos/raise → backend
```

### Local Development

For local testing with Vercel:

```bash
# Terminal 1: Start backend on port 8080
cd backend
go run .

# Terminal 2: Start frontend (Vite dev server)
cd frontend
npm run dev  # Runs on http://localhost:5173

# Vite dev server proxies /api requests to http://localhost:8080
# See frontend/vite.config.js for proxy config
```

## Deployment Instructions

### 1. Push to GitHub

Ensure all changes are committed:

```bash
git add -A
git commit -m "config: set up Vercel services deployment"
git push origin main
```

### 2. Import on Vercel

1. Go to **https://vercel.com/dashboard**
2. Click **"Add New"** → **"Project"**
3. Select your GitHub repository: `Rishabh2815/rakshagrid`
4. Vercel auto-detects `vercel.json` and configures both services automatically
5. Click **"Deploy"**

### 3. Environment Variables (if needed)

In Vercel Dashboard → Project Settings → Environment Variables:

```
# Backend (Go)
POSTGRES_DSN=postgres://...
REDIS_ADDR=redis.railway.app:6379
REDIS_PASSWORD=***
JWT_SECRET=your-production-secret
SOS_ACK_SLA=60s
DISPATCH_RADIUS_KM=2.0

# Frontend (React)
# No environment variables needed — uses /api hardcoded
```

**Note:** If you still want to use Railway for Postgres/Redis:
1. Deploy on Railway (as before) and get credentials
2. Set those credentials as Vercel env vars
3. They'll be injected at build/runtime

### 4. Verify Deployment

```bash
# Frontend loads
curl https://rakshagrid.vercel.app/

# Backend responds
curl https://rakshagrid.vercel.app/api/health
# Expected: {"status":"ok"}

# SOS endpoint exists
curl -X POST https://rakshagrid.vercel.app/api/sos/raise \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer dummy-token" \
  -d '{"latitude":0,"longitude":0}'
```

## Comparison: Previous Setup vs. Vercel Services

| Aspect | Railway Setup | Vercel Services |
|--------|---------------|-----------------|
| **Frontend** | Vercel | Vercel (integrated) |
| **Backend** | Railway | Vercel (integrated) |
| **Domain** | Two domains (vercel.app + railway.app) | One domain (vercel.app) |
| **CORS** | Required | Not needed (same origin) |
| **Cost** | $5–10/month (Railway backend) | Free tier available |
| **Databases** | Railway PostgreSQL/Redis | Must use external services |
| **Scaling** | Separate scaling per service | Vercel auto-scaling |

## Using External Databases

Vercel doesn't include managed databases. For Postgres/Redis:

### Option A: Railway (Recommended)

1. Create Railway project with Postgres + Redis (free tier)
2. Get connection strings from Railway
3. Set as Vercel env vars
4. Backend uses them at runtime

### Option B: Other Services

- **Neon** (serverless Postgres)
- **Upstash** (serverless Redis)
- **AWS RDS** (traditional RDS)
- **DigitalOcean** (managed databases)

## Local Development with `vercel dev`

To test Vercel services locally (before pushing):

```bash
# Requires Vercel CLI
npm install -g vercel

# Run all services together
vercel dev

# Vercel injects bindings and environment variables
# Frontend: http://localhost:3000
# Backend: accessible at /api from frontend
```

## Troubleshooting

### Frontend can't reach backend
- ❌ Check CORS headers — not needed in Vercel services
- ✅ Verify rewrite in `vercel.json`: `/api/(.*)` → `backend`
- ✅ Check backend logs in Vercel dashboard

### WebSocket not connecting
- Ensure `wss://` is used (HTTPS required in production)
- Verify `/api/stream` endpoint exists in backend
- Check browser DevTools → Network → WS tab

### Build fails
- Check Vercel logs → select failed service (frontend or backend)
- Frontend: look for Node.js/npm errors
- Backend: look for Go compilation errors

## Next Steps

1. **Deploy on Vercel** using steps above
2. **Implement backend logic:**
   - SOS dispatch orchestration
   - Redis geospatial matching
   - WebSocket real-time updates
   - Incident audit logging
3. **Add database migrations** (if using external Postgres)
4. **Test end-to-end** SOS flow

## References

- **Vercel Services Docs:** https://vercel.com/docs/services
- **Vercel CLI:** https://vercel.com/cli
- **System Design:** See `RakshaGrid_System_Design.md`
