<h1 align="center">Ask Me Anything (Client + Server)</h1>

<!-- markdownlint-disable MD033 MD041 -->
<div align="center">
  <p>Real-time Q&A application unified in a single repository with a Next.js frontend and a Go backend.</p>
  
  <p>
    <img alt="Go" src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white" />
    <img alt="chi" src="https://img.shields.io/badge/chi-5.2.2-4BA0F6" />
    <img alt="gorilla/websocket" src="https://img.shields.io/badge/gorilla/websocket-1.5.3-7C4DFF" />
    <img alt="pgx" src="https://img.shields.io/badge/pgx-5.7.5-2E8B57" />
    <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-latest-336791?logo=postgresql&logoColor=white" />
    <img alt="Next.js" src="https://img.shields.io/badge/Next.js--000000?logo=nextdotjs&logoColor=white" />
    <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript&logoColor=white" />
    <img alt="Docker Compose" src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" />
  </p>
</div>

---

## Overview

This monorepo contains:

- `client/` — Frontend in Next.js (TypeScript)
- `server/` — API in Go (WebSocket + REST) with PostgreSQL

Real-time communication happens via WebSockets (`/subscribe/{room_id}`); the rest is served over REST under `/api`.

## Structure

```text
ask-me-anything/
├── client/                     # Frontend Next.js (TS)
│   ├── app/                    # Routes and pages (App Router)
│   ├── next.config.ts          # Next.js configuration
│   ├── eslint.config.mjs       # Frontend lint
│   └── ...                     
└── server/                     # Backend in Go
  ├── cmd/wsrs/               # Server entrypoint
  ├── internal/               # API, auth, store (pgx/sqlc), websockets, etc.
  ├── docker-compose.yml      # Local stack (API + DB + pgAdmin)
  ├── Dockerfile              # API build
  ├── Makefile                # Dev tasks
  └── docs/                   # API documentation
```

## Tech stack

- Backend: Go 1.25, Chi, Gorilla WebSocket, pgx, SQLC, PostgreSQL
- Frontend: Next.js (TS)
- Infra: Docker & Docker Compose

## Quick start

### Option A — Docker Compose (recommended, API/DB only)

```bash
cd server
# bring up API + PostgreSQL (+ pgAdmin)
docker compose up -d

# URLs
# API:       http://localhost:8080
# Health:    http://localhost:8080/health
# WebSocket: ws://localhost:8080/subscribe/{room_id}
```

The frontend runs in parallel in the `client/` directory.

### Option B — Run separately

- Backend (Go):

```bash
cd server
make migrate-up   # if needed
make run          # or: go run ./cmd/wsrs
```

- Frontend (Next.js):

```bash
cd client
# install dependencies with your preferred package manager and start
# npm install
# npm run dev
```

## Environment variables (API)

Create `server/.env` (or export env vars):

```env
LOG_LEVEL=info
WSRS_DATABASE_HOST=localhost
WSRS_DATABASE_PORT=5432
WSRS_DATABASE_USER=postgres
WSRS_DATABASE_PASSWORD=password
WSRS_DATABASE_NAME=wsrs_db
```

Tip: the `server/Makefile` helps with migrations, logs, and Docker rebuilds.

## Detailed documentation

- API: see `server/README.md` and `server/docs/`
- Frontend: configure `.env.local` as needed and run the dev server

## Monorepo notes

- Unified repository to simplify development and integrated versioning.
- `client/` and `server/` keep responsibilities isolated but share the same lifecycle (PRs/releases).
- The `server/` Compose provides DB + API; the client runs separately.

---

Made with ❤️.
