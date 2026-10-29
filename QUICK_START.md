# 🚀 RakshaGrid: Deploy to Railway + Vercel in 5 Minutes

Your project is **ready to deploy**. Here's the quickest path to production.

## Step 1: Deploy Backend to Railway (3 min)

```
1. Visit https://railway.app
2. Click "New Project" → "Deploy from GitHub repo"
3. Select "Rishabh2815/rakshagrid"
4. Click "+ Add Service" → PostgreSQL
5. Click "+ Add Service" → Redis
6. Click "Deploy" 
7. Wait for build to complete (~2 min)
8. Note your backend URL: https://rakshagrid-api.railway.app (example)
```

✅ Backend is now live.

## Step 2: Deploy Frontend to Vercel (2 min)

```
1. Visit https://vercel.com
2. Click "Add New Project" → Select "rakshagrid"
3. Build Command: cd frontend && npm install && npm run build
4. Output Directory: frontend/dist
5. Add env var:
   VITE_BACKEND_URL=https://rakshagrid-api.railway.app
6. Click "Deploy"
7. Wait for build (~1 min)
8. Note your frontend URL: https://rakshagrid-xyzabc.vercel.app
```

✅ Frontend is now live.

## Step 3: Update CORS (1 min)

Edit `backend/internal/api/router.go` line ~15:

```go
AllowedOrigins: []string{
  "http://localhost:3000",
  "https://rakshagrid-xyzabc.vercel.app",  // ← ADD YOUR VERCEL URL HERE
}
```

Then:
```bash
git add backend/internal/api/router.go
git commit -m "chore: update CORS for production domain"
git push origin main
```

Railway auto-redeploys. (Wait 1 min.)

## ✅ Done!

Open your Vercel frontend URL. You should see:
- Login page
- SOS button (when logged in)
- ✅ No errors in browser console

---

## 🔗 Your Live URLs

- **Frontend:** https://rakshagrid-xyzabc.vercel.app
- **Backend:** https://rakshagrid-api.railway.app
- **Health Check:** https://rakshagrid-api.railway.app/health

---

## 🧪 Test the SOS Flow

1. Click Login → use demo credentials (shown on page)
2. You're redirected to SOS button
3. Click the big red button
4. Browser DevTools → Network tab → see `/v1/sos/raise` POST request
5. Response should be `{"status":"sos_received"}`

---

## ⚠️ If Something Breaks

### Backend won't start
```bash
# Check logs in Railway Dashboard
# Backend Service → Logs tab
# Look for "connection failed" or "bind: address already in use"
```

### Frontend build fails
```bash
# Check logs in Vercel Dashboard
# Deployments → Click latest → Runtime Logs
# Look for "vite build" errors
```

### CORS errors
```bash
# Browser console shows: "Access to XMLHttpRequest blocked by CORS policy"
# Solution: Add your Vercel URL to AllowedOrigins (Step 3 above)
```

---

## 📖 Full Docs

- **Detailed Setup:** `RAILWAY_DEPLOYMENT.md`
- **Complete Setup Info:** `DEPLOYMENT_READY.md`
- **Architecture:** `RakshaGrid_System_Design.md`

---

**You're live!** 🎉 Next: complete the backend implementation (dispatch, responders, audit logging).
