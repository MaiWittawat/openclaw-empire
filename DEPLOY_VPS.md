# OpenClaw Empire VPS Deploy

This backend is designed to run on the same Ubuntu VPS as OpenClaw Gateway and connect to the shared PostgreSQL infra container.

## Prerequisites

- OpenClaw Gateway is already running on the host
- PostgreSQL is already running in the shared Docker infra on `web-proxy-net`
- `OPENCLAW_GATEWAY_TOKEN` is ready

## Environment

Copy `.env.example` to `.env` and fill at least:

```env
DB_HOST=psql
DB_PORT=5432
DB_USER=claw_admin
DB_PASS=claw_admin
DB_NAME=claw_admin

OPENCLAW_GATEWAY_WS_URL=ws://host.docker.internal:18789/ws
OPENCLAW_GATEWAY_TOKEN=your_gateway_token
OPENCLAW_DEFAULT_SESSION_KEY=main
BACKEND_URL=http://claw-empire-backend:8080
```

`OPENCLAW_DEVICE_STORE_PATH` should stay on a persistent volume so the device token can be reused after restarts.

## Run

```bash
docker compose up -d --build
```

## Notes

- `backend` joins external network `web-proxy-net` to reach the shared PostgreSQL container `psql`
- `backend` reaches OpenClaw Gateway on the VPS host through `host.docker.internal`
- Task history and event logs are stored in PostgreSQL
- Realtime agent updates come from OpenClaw Gateway WebSocket events

## API

- `GET /api/agents`
- `GET /api/tasks`
- `GET /api/tasks/:id/events`
- `POST /api/tasks`
- `GET /api/stats`
- `GET /ws`
