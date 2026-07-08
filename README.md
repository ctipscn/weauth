# WeAuth

WeAuth is a proof-of-work based verification and anti-abuse platform. It provides a captcha-like user experience without image recognition or slider puzzles: the browser performs a small SHA-256 proof-of-work challenge and the backend verifies the result before issuing a short-lived token.

**Live demo:** https://www.weauth.cn

## Features

- Proof-of-work human verification with adjustable difficulty.
- Embeddable browser widget (`weauth.js`) for third-party sites.
- Site management, secrets, verification statistics, and IP policies.
- OIDC login support for an external identity provider.
- Go backend built with Gin, GORM, PostgreSQL, and JWT.
- React/Vite frontend with dashboard and documentation pages.

## Repository Layout

```text
backend/    Go API service, database models, migrations, PoW verification
frontend/   React dashboard and embeddable widget build
```

## Quick Start

### Backend

```bash
cd backend
cp .env.example .env
go mod download
go run main.go
```

Configure PostgreSQL and OIDC values in `backend/.env` before using the login flow.

### Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Build the dashboard and widget:

```bash
npm run build
```

## Important Configuration

- `JWT_SECRET` must be a strong random value in production.
- `OIDC_CLIENT_SECRET` must never be committed.
- `CORS_ORIGINS` should list only trusted frontend origins.
- The public widget can be configured with `VITE_WIDGET_BASE_URL` and `VITE_API_BASE_URL`.

## License

Apache-2.0
