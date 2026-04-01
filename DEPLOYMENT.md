# Deployment Guide

This project is prepared for:
- Frontend: Vercel
- Backend: Railway
- Database: Neon (PostgreSQL)

## 1. Neon Database

1. Create a Neon project and database.
2. Copy the connection string from Neon (prefer pooled connection if available).
3. Ensure the URL includes SSL mode, for example: `...?sslmode=require`.
4. Save this value for Railway as `DATABASE_URL`.

## 2. Backend on Railway

1. Create a new Railway service from this repository.
2. Set the service root directory to `backend` (if using monorepo settings).
3. Railway will build from `backend/Dockerfile`.
4. Add environment variables:

```env
DATABASE_URL=<your_neon_connection_string>
DB_SSLMODE=require
FRONTEND_ORIGINS=https://<your-vercel-domain>
JWT_SECRET=<strong_random_secret>
JWT_EXPIRY_HOURS=24
ADMIN_NAME=Administrator
ADMIN_USERNAME=admin
ADMIN_PASSWORD=<strong_password>
UPLOAD_DIR=/app/uploads
```

Notes:
- `PORT` is injected by Railway automatically and now supported by backend.
- Health check endpoint: `/api/v1/health`
- You can allow previews with wildcard, for example:
  `FRONTEND_ORIGINS=https://<your-app>.vercel.app,https://*.vercel.app`

## 3. Frontend on Vercel

1. Create a Vercel project from this repository.
2. Set project root directory to `frontend`.
3. Add environment variable:

```env
NEXT_PUBLIC_API_URL=https://<your-railway-domain>/api/v1
```

4. Deploy and copy the generated Vercel domain.

## 4. Final Wiring

1. Update Railway `FRONTEND_ORIGINS` with your actual Vercel domain.
2. Redeploy Railway backend after changing env vars.
3. Test login and protected APIs from Vercel frontend.

## 5. Local Environment Templates

- Backend template: `backend/.env.example`
- Frontend template: `frontend/.env.example`
