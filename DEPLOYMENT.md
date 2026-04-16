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
JWT_SECRET=<strong_random_secret_at_least_32_chars>
JWT_EXPIRY_HOURS=24
ADMIN_NAME=Administrator
ADMIN_USERNAME=admin
ADMIN_PASSWORD=<strong_password>
UPLOAD_DIR=/app/uploads

# Database logging: silent, error, warn, info (default: warn)
DB_LOG_LEVEL=warn

# Set to true to keep uploaded Excel files after import confirmation (default: false)
IMPORT_KEEP_FILES=false
```

### Secret Policy

- **`JWT_SECRET`**: Must be at least 32 characters and must not be a known default/placeholder value. The application will refuse to start if this requirement is not met.
- **`ADMIN_PASSWORD`**: Must be a strong, non-default password. The application will refuse to bootstrap the admin user with weak or placeholder passwords.
- **`DB_PASSWORD`**: Must be set in the environment; no hardcoded defaults are used.

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

> **Note**: The frontend uses a local font file (`GeistVF.woff`) and does not require internet access to Google Fonts during build.

## 4. Final Wiring

1. Update Railway `FRONTEND_ORIGINS` with your actual Vercel domain.
2. Redeploy Railway backend after changing env vars.
3. Test login and protected APIs from Vercel frontend.

## 5. Local Environment Templates

- Backend template: `backend/.env.example`
- Frontend template: `frontend/.env.example`

### New Environment Variables Reference

| Variable | Default | Description |
|---|---|---|
| `DB_LOG_LEVEL` | `warn` | GORM log level: `silent`, `error`, `warn`, `info` |
| `IMPORT_KEEP_FILES` | `false` | Keep uploaded Excel files after successful import |
| `JWT_SECRET` | _(required)_ | Must be ≥ 32 chars, non-default |
| `ADMIN_PASSWORD` | _(required)_ | Must be strong, non-default |
