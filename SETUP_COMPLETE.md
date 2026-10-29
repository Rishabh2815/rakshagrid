# Railway Deployment Checklist for RakshaGrid

## ✅ Setup Complete

Your project has been reorganized for Railway deployment with the following structure:

```
rakshagrid/
├── backend/
│   ├── main.go
│   ├── go.mod / go.sum
│   ├── Dockerfile           (multi-stage Go build)
│   ├── Procfile
│   └── internal/
│       ├── config/          (environment configuration)
│       ├── db/              (Postgres + Redis connections)
│       ├── auth/            (JWT + password hashing)
│       └── api/             (HTTP router + basic handlers)
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js       (with VITE_BACKEND_URL support)
│   ├── Dockerfile           (Node.js + Nginx)
│   ├── nginx.conf
│   └── [your React files]
├── railway.json             (Railway build config)
└── RAILWAY_DEPLOYMENT.md    (this guide)
```

## 🚀 Next Steps (Follow RAILWAY_DEPLOYMENT.md)

1. **Push to GitHub** (all changes committed)
2. **Create Railway project** and link to your GitHub repo
3. **Add PostgreSQL + Redis** services (auto-managed)
4. **Deploy backend** on Railway (URL: `https://rakshagrid-api.railway.app`)
5. **Deploy frontend** on Vercel with `VITE_BACKEND_URL` env var pointing to Railway
6. **Verify** health check and CORS configuration

## 📋 What Was Created

### Backend Scaffolding

- ✅ **Dockerfile** - Multi-stage Go build optimized for Railway
- ✅ **Procfile** - Defines `web: ./server` entrypoint
- ✅ **internal/config** - Reads env vars for DB, Redis, JWT, CORS
- ✅ **internal/db** - PostgreSQL + Redis connection setup
- ✅ **internal/auth** - JWT token generation + password hashing (MVP)
- ✅ **internal/api** - HTTP router with CORS, basic health check

### Frontend Updates

- ✅ **vite.config.js** - Production-ready with `VITE_BACKEND_URL` support
- ✅ **Dockerfile** - Multi-stage Node.js + Nginx build
- ✅ **nginx.conf** - Already in place for static asset serving

### Configuration

- ✅ **railway.json** - Tells Railway to build `backend/Dockerfile`
- ✅ **RAILWAY_DEPLOYMENT.md** - Complete step-by-step guide

## ⚠️ Important Notes

### Backend is MVP Scaffolding

The backend implements:
- ✅ Config loading from environment
- ✅ Postgres + Redis connection setup
- ✅ JWT auth service (SHA256 hashing — use bcrypt in production)
- ✅ Basic HTTP router with CORS middleware
- ✅ Health check endpoint (`/health`)
- ✅ Skeleton SOS + responder handlers

**Missing (you need to implement):**
- 🔴 Core SOS dispatch logic
- 🔴 Geospatial responder matching (Redis GEOSEARCH)
- 🔴 WebSocket real-time updates
- 🔴 Audit logging
- 🔴 Database schema + migrations

See `RakshaGrid_System_Design.md` for the full architectural requirements.

### Environment Variables

Railway will auto-inject these from linked services:

```
POSTGRES_DSN=${{ Postgres.DATABASE_URL }}
REDIS_ADDR=redis:6379  # Railway internal DNS
REDIS_PASSWORD=${{ Redis.REDIS_PASSWORD }}
```

You must manually set:

```
HTTP_PORT=8080
JWT_SECRET=your-secret-key-here
SOS_ACK_SLA=60s
DISPATCH_RADIUS_KM=2.0
FRONTEND_URL=https://your-vercel-project.vercel.app
```

### CORS Settings

Currently allows `https://*` for development convenience. Before production, update `internal/api/router.go`:

```go
AllowedOrigins: []string{
  "https://your-vercel-project.vercel.app",
  // Remove "https://*"
}
```

## 🔗 Quick Links

- **Railway Dashboard:** https://railway.app
- **Vercel Dashboard:** https://vercel.com
- **Full Deployment Guide:** See `RAILWAY_DEPLOYMENT.md` in this directory
- **System Design:** See `RakshaGrid_System_Design.md` for architecture reference

## 🐛 Troubleshooting

### Backend won't compile
```bash
cd backend
go mod tidy
go build -o server .
```

### Dockerfile fails to build
- Ensure `backend/Dockerfile` exists and is readable
- Verify `backend/main.go` and `go.mod` are in the backend directory
- Check file permissions: `chmod 644 backend/Dockerfile`

### Missing dependencies
```bash
cd backend
go mod download
```

### Test locally before deploying to Railway
```bash
cd backend
HTTP_PORT=8080 POSTGRES_DSN="..." REDIS_ADDR="..." JWT_SECRET="dev" ./server
```

---

**All set!** Follow the deployment steps in `RAILWAY_DEPLOYMENT.md` to go live.
