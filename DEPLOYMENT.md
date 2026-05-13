# Deployment Guide

This project is prepared for:

- Frontend: Vercel
- Backend: Railway
- Database: Neon PostgreSQL

## VPS Deployment (Docker: Backend + PostgreSQL)

Use this path when deploying the backend service directly to your own cloud VPS.

### Files used

- `docker-compose.vps.yml`
- `.env.vps` (copy from `.env.vps.example`)

### 1. Prepare VPS

Install Docker Engine and Docker Compose plugin on the VPS, then clone this repository.

### 2. Create production env file

```bash
cp .env.vps.example .env.vps
```

Then update at least:

- `DB_PASSWORD`
- `JWT_SECRET`
- `ADMIN_PASSWORD`
- `FRONTEND_ORIGINS`

### 3. Start backend + database

```bash
docker compose --env-file .env.vps -f docker-compose.vps.yml up -d --build
```

### 4. Verify health

```bash
docker compose --env-file .env.vps -f docker-compose.vps.yml ps
curl http://127.0.0.1:8080/api/v1/health
```

### 5. Persisted data

Compose uses named volumes:

- `pgdata` for PostgreSQL data
- `uploads` for uploaded/imported files

Do not remove volumes unless you intentionally want to wipe data.

### 6. Update or redeploy

```bash
git pull
docker compose --env-file .env.vps -f docker-compose.vps.yml up -d --build
```

### 7. Recommended production fronting

- Put Nginx/Caddy/Traefik in front of `:8080` for TLS termination.
- Keep `COOKIE_SECURE=true` and `COOKIE_SAMESITE=None` for HTTPS cross-site cookie auth.
- Keep `FRONTEND_ORIGINS` as an exact allowlist (no wildcard).

## 1. Neon Database

1. Create a Neon project and database.
2. Copy the connection string from Neon. Prefer pooled connection if available.
3. Ensure the URL includes SSL mode, for example: `...?sslmode=require`.
4. Save this value for Railway as `DATABASE_URL`.

## 2. Backend on Railway

1. Create a new Railway service from this repository.
2. Set the service root directory to `backend` if using monorepo settings.
3. Railway will build from `backend/Dockerfile`.
4. Add environment variables:

```env
DATABASE_URL=<your_neon_connection_string>
DB_SSLMODE=require
FRONTEND_ORIGINS=https://<your-production-vercel-domain>
JWT_SECRET=<strong_random_secret_at_least_32_chars>
JWT_EXPIRY_HOURS=24
ADMIN_NAME=Administrator
ADMIN_USERNAME=admin
ADMIN_PASSWORD=<strong_password_at_least_12_chars_with_letters_and_numbers>
UPLOAD_DIR=/app/uploads
GIN_MODE=release
TRUSTED_PROXIES=
MAX_BODY_BYTES=10485760
IMPORT_MAX_BYTES=57671680
IMPORT_JSON_BYTES=20971520

# Database logging: silent, error, warn, info (default: warn)
DB_LOG_LEVEL=warn

# Set to true to keep uploaded Excel files after import confirmation (default: false)
IMPORT_KEEP_FILES=false
```

### Secret Policy

- `JWT_SECRET`: Must be at least 32 characters and must not be a known default or placeholder value. The application refuses to start when this requirement is not met.
- `ADMIN_PASSWORD`: Must be at least 12 characters, include letters and numbers, and must not be a known placeholder. The application refuses to bootstrap or sync the admin user with weak passwords.
- `DB_PASSWORD`: Must be set when using discrete DB variables instead of `DATABASE_URL`; no hardcoded database password default is used.

Notes:

- `PORT` is injected by Railway automatically and is supported by the backend.
- Health check endpoint: `/api/v1/health`
- Cookie-authenticated production deployments must use exact origins in `FRONTEND_ORIGINS`. Do not use wildcard Vercel origins for production cookie auth.
- If preview deployments need API access, add each trusted preview origin explicitly and remove it when it is no longer needed.

## 3. Frontend on Vercel

1. Create a Vercel project from this repository.
2. Set project root directory to `frontend`.
3. Add environment variable:

```env
NEXT_PUBLIC_API_URL=https://<your-railway-domain>/api/v1
```

4. Deploy and copy the generated Vercel domain.

The frontend uses a local font file and does not require internet access to Google Fonts during build.

## 4. Final Wiring

1. Update Railway `FRONTEND_ORIGINS` with your actual Vercel domain.
2. Redeploy Railway backend after changing env vars.
3. Test login and protected APIs from the Vercel frontend.

## 5. Local Environment Templates

- Root Compose template: `.env.example`
- Backend template: `backend/.env.example`
- Frontend template: `frontend/.env.example`

### Environment Variables Reference

| Variable | Default | Description |
|---|---|---|
| `DB_LOG_LEVEL` | `warn` | GORM log level: `silent`, `error`, `warn`, `info` |
| `IMPORT_KEEP_FILES` | `false` | Keep uploaded Excel files after successful import |
| `JWT_SECRET` | required | Must be >= 32 chars, non-default |
| `ADMIN_PASSWORD` | required | Must be >= 12 chars with letters and numbers, non-default |
| `GIN_MODE` | `debug` | Set to `release` in production |
| `TRUSTED_PROXIES` | empty | Comma-separated trusted proxy CIDRs/IPs; empty disables trusted proxies |
| `MAX_BODY_BYTES` | `10485760` | Default API request body limit in bytes |
| `IMPORT_MAX_BYTES` | `57671680` | Excel upload request limit in bytes |
| `IMPORT_JSON_BYTES` | `20971520` | Import confirmation JSON request limit in bytes |

## 6. Local Docker Compose

The root Compose stack runs PostgreSQL, the Go API, and the Next.js frontend.

1. Create a root `.env` from `.env.example`.
2. Replace the placeholder secrets:
   - `DB_PASSWORD`
   - `JWT_SECRET`
   - `ADMIN_PASSWORD`
3. Start the stack:

```bash
docker compose up --build
```

Local service URLs:

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080/api/v1`
- Backend health: `http://localhost:8080/api/v1/health`

Compose intentionally fails fast when required secrets are missing. The backend still performs its own JWT/admin-password validation during startup.

## 7. Local Without Docker

1. Start PostgreSQL locally.
2. Create `backend/.env` from `backend/.env.example`.
3. Replace `JWT_SECRET` and `ADMIN_PASSWORD` with strong non-default values.
4. Confirm `FRONTEND_ORIGINS` includes the frontend origin, normally `http://localhost:3000`.
5. Start the backend:

```bash
cd backend
go run ./cmd/server
```

6. Create `frontend/.env` from `frontend/.env.example`.
7. Confirm `NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1`.
8. Start the frontend:

```bash
cd frontend
npm install
npm run dev
```

For cross-origin cookie auth, keep these values aligned:

- Frontend requests must use credentials through the shared API client.
- Backend `FRONTEND_ORIGINS` must list the exact frontend origin.
- Unsafe cookie-authenticated mutations are rejected unless `Origin` or `Referer` matches an exact configured frontend origin.
- For HTTPS cross-site production, use `COOKIE_SAMESITE=None` and `COOKIE_SECURE=true`.
- For local HTTP development, `COOKIE_SAMESITE=Lax` and `COOKIE_SECURE=false` are expected.
