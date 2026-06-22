# 🏊 aqua-os-backend

REST API for water polo club management — exercises, levels, and training data.

## Architecture

Go service using chi router with SQLite storage. Serves as the core data API for the aqua-os frontend and AI crews.

```
cmd/server/    → entrypoint
internal/      → handlers, models, database
data/          → SQLite database file
```

## Prerequisites

- Go 1.21+
- GCC (for SQLite CGO bindings)

## Run Locally

```bash
make run
# or
go run ./cmd/server
```

Server starts on **http://localhost:8080**

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `JWT_SECRET` | JWT signing key | — |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/levels` | List all levels |
| GET | `/api/levels/{id}` | Get level by ID |
| POST | `/api/exercises` | Create exercise |
| GET | `/api/exercises` | List exercises |
| DELETE | `/api/exercises/{id}` | Delete exercise |

## Docker

```bash
docker build -t aqua-os-backend .
docker run -p 8080:8080 aqua-os-backend
```

## Related Repos

| Repo | Description |
|------|-------------|
| [aqua-os-web](../aqua-os-web) | React frontend |
| [aqua-os-crew](../aqua-os-crew) | AI agents (CrewAI) |
| [aqua-os-calendar](../aqua-os-calendar) | Game calendar service |
| [aqua-os-infrastructure](../aqua-os-infrastructure) | Terraform AWS infra |

## License

GPL-3.0
