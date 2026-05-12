# Dashboard Procurement Frontend

Next.js App Router frontend for the procurement dashboard.

## Stack

- Next.js, React, and TypeScript
- Refine resource routing
- Ant Design components and tokens
- Axios API client with credentials enabled
- Recharts for overview charts

## Local Development

Create `frontend/.env` from `frontend/.env.example`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

Install dependencies and run the app:

```bash
npm install
npm run dev
```

The app runs at `http://localhost:3000`.

## Production Build

```bash
npm run build
npm run start
```

On this Windows workspace, use `cmd /c npm run ...` when PowerShell blocks npm shims.

## Docker

The root `docker-compose.yml` builds this frontend from `frontend/Dockerfile`. The API URL is baked into the production bundle through the `NEXT_PUBLIC_API_URL` build argument, so set it in the root `.env` before running Compose.
