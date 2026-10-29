# RakshaGrid Railway Deployment Guide

This guide walks you through deploying RakshaGrid on Railway with PostgreSQL, Redis, and a Go backend, then connecting a Vercel-hosted frontend.

## Prerequisites

- GitHub account with your rakshagrid repo pushed
- Railway account (free tier: https://railway.app)
- Vercel account (free tier: https://vercel.com)
- Your repo pushed to GitHub with all changes

## Architecture

```
┌─────────────────────┐
│   Vercel Frontend   │
│   (React + Vite)    │
└──────────┬──────────┘
           │ /v1 API calls
           ▼
┌─────────────────────────────┐
│   Railway Backend Service   │
│   (Go + WebSocket)          │
└──────────┬──────────────────┘
           │
    ┌──────┴───────┐
    ▼              ▼
  PostgreSQL     Redis
  (Railway)      (Railway)
```

## Step 1: Prepare Your Repository

1. **Push all changes to GitHub:**
   ```bash
   cd C:\Users\Rishabh\CodeRepo\rakshagrid
   git add -A
   git commit -m "refactor: reorganize into frontend/backend structure for Railway deployment"
   git push origin main
   ```

2. **Verify the structure:**
   - ✅ `backend/main.go`, `go.mod`, `go.sum`, `Dockerfile`, `Procfile`
   - ✅ `backend/internal/` (config, db, auth, api)
   - ✅ `frontend/package.json`, `vite.config.js`, `index.html`, `Dockerfile`, `nginx.conf`
   - ✅ `railway.json` (root)

## Step 2: Deploy Backend to Railway

### 2.1 Create Railway Project

1. Go to https://railway.app and **Sign in with GitHub**
2. Click **"New Project"** → **"Deploy from GitHub repo"**
3. Search for and select `Rishabh2815/rakshagrid`
4. Railway auto-detects the `railway.json` and `backend/Dockerfile`

### 2.2 Add PostgreSQL Database

1. In your Railway project, click **"+ Add Service"**
2. Select **"PostgreSQL"** from the available plugins
3. Railway auto-generates:
   - `POSTGRES_DSN` (connection string)
   - Database name: `railway`
   - User: `postgres`
   - Password: (auto-generated)

### 2.3 Add Redis Cache

1. Click **"+ Add Service"** again
2. Select **"Redis"** from plugins
3. Railway auto-generates:
   - `REDIS_ADDR` (e.g., `redis.railway.internal:6379`)
   - `REDIS_PASSWORD` (if set)

### 2.4 Configure Backend Service Environment

Click on the **"backend"** service (the one built from your `Dockerfile`), then **"Variables"** tab:

Set these environment variables:

```
HTTP_PORT=8080
POSTGRES_DSN=${{ Postgres.DATABASE_URL }}
REDIS_ADDR=redis:6379
REDIS_PASSWORD=${{ Redis.REDIS_PASSWORD }}
JWT_SECRET=your-production-secret-here-change-me
SOS_ACK_SLA=60s
DISPATCH_RADIUS_KM=2.0
FRONTEND_URL=https://your-vercel-project.vercel.app
```

**Note:** Railway uses `${{ Service.VARIABLE }}` syntax to inject values from linked services.

### 2.5 Deploy

Click **"Deploy"**. Railway will:
1. Build from `backend/Dockerfile`
2. Start the Go server on port 8080
3. Assign you a public URL (e.g., `https://rakshagrid-api.railway.app`)

**Watch the logs** to ensure the server starts without errors:
- Should see: `rakshagrid server listening on :8080`
- Should NOT see database connection failures

## Step 3: Deploy Frontend to Vercel

### 3.1 Push to GitHub (already done in Step 1)

### 3.2 Import to Vercel

1. Go to https://vercel.com and **Sign in with GitHub**
2. Click **"Add New..."** → **"Project"**
3. Search for `rakshagrid` and click **"Import"**
4. Vercel auto-detects it as a Vite project

### 3.3 Configure Environment Variables

In the **Environment Variables** section, add:

```
VITE_BACKEND_URL=https://rakshagrid-api.railway.app
```

(Replace with your actual Railway backend URL from Step 2.5)

### 3.4 Deploy

Click **"Deploy"**. Vercel will:
1. Build the React app with `npm run build`
2. Deploy to `https://<your-project>.vercel.app`
3. All `/v1/*` requests proxy to your Railway backend

## Step 4: Verify Everything Works

### 4.1 Test Backend Health

```bash
curl https://rakshagrid-api.railway.app/health
# Expected: {"status":"ok"}
```

### 4.2 Test Frontend

1. Open https://your-vercel-project.vercel.app in your browser
2. Open browser DevTools → **Network** tab
3. Trigger an SOS (or navigate to an API endpoint)
4. You should see `/v1/*` requests going to your Railway backend ✅

### 4.3 Check Logs

**Railway Backend Logs:**
- Railway Dashboard → Backend Service → Logs tab
- Watch for connection confirmations: `postgres connection succeeded`, `redis connection succeeded`

**Vercel Frontend Logs:**
- Vercel Dashboard → Deployments → Click latest deployment → Runtime Logs

## Step 5: Update CORS (Backend)

If you encounter CORS errors, update `internal/api/router.go` to include your Vercel domain:

```go
AllowedOrigins: []string{
  "http://localhost:3000",
  "http://localhost:5173",
  "https://your-vercel-project.vercel.app",  // ← Add this
  "https://*",  // Allow all HTTPS in production (optional)
}
```

Then push and redeploy on Railway.

## Step 6: Database Migrations (Optional)

If you have SQL schema files in `backend/migrations/`:

1. Connect to Railway PostgreSQL:
   ```bash
   psql $POSTGRES_DSN < backend/migrations/001_init.sql
   ```

2. Or create a migration runner in your Go code (recommended)

## Troubleshooting

### Backend won't start
- **Check logs:** Railway Dashboard → Backend Logs
- **Common issues:**
  - Missing `POSTGRES_DSN` environment variable
  - Redis connection failed
  - Port already in use (Railway handles this)

### Frontend can't reach backend
- **Check proxy config:** `frontend/vite.config.js` has correct `VITE_BACKEND_URL`
- **Check CORS:** Backend must allow your Vercel domain
- **Check Network tab** in browser DevTools for the actual API call

### Database connection errors
- Verify `POSTGRES_DSN` is set and correct
- Check if database service is running in Railway dashboard
- Try manually connecting via psql with the DSN

## Scaling & Next Steps

1. **Add monitoring:** Railway → Project → Monitoring tab (Grafana integration available)
2. **Enable auto-deploy on push:** Railway → GitHub integration (auto-deploy on `git push`)
3. **Database backups:** Railway PostgreSQL plugin offers automated backups
4. **Custom domain:** Vercel / Railway both support custom domains (paid)

## Cost Estimate

- **Railway:** ~$5–10/month (backend + Postgres + Redis)
- **Vercel:** Free tier (frontend)
- **Total:** ~$5–10/month

---

**Deployment complete!** Your RakshaGrid is now live. Next: complete the backend implementation in `internal/dispatch`, `internal/notify`, `internal/audit` packages to enable full SOS functionality.
