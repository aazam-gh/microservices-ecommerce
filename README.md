# E-Commerce Platform

Microservices-based backend platform with a running local stack (Docker Compose), centralized auth through Keycloak, and gateway routing through Nginx.

This README reflects the current repository state as of March 31, 2026.

## Current Status

Implemented services:
- `user-service` (Go + Fiber + PostgreSQL + Redis + Keycloak introspection)
- `search-service` (Go + Fiber + Elasticsearch + optional Redis cache)

Scaffolded but not implemented yet:
- `order-service`
- `payment-service`
- `notification-service`

Work in progress:
- `product-service` (.NET 8 + ASP.NET Core + PostgreSQL)

## Tech Stack

- Go `1.25.x` (user/search services)
- .NET `8.0` (product service)
- PostgreSQL 16
- Redis 7
- Kafka + Zookeeper + Kafdrop
- Elasticsearch 8 + Kibana
- Keycloak 24
- Nginx
- Prometheus + Grafana + Jaeger

## Repository Layout

- `services/user-service` - user profile, address management, token validation endpoint
- `services/product-service` - product + category service (WIP)
- `services/search-service` - Elasticsearch query API (`/search?q=...`)
- `services/order-service` - scaffold (`go.mod` only)
- `services/payment-service` - scaffold (`go.mod` only)
- `services/notification-service` - scaffold (`go.mod` only)
- `infra` - Docker Compose, Nginx, Keycloak realm export, Prometheus/Grafana, Terraform, K8s manifests
- `docs/openapi` - placeholder only (`.gitkeep`)

## Prerequisites

- Docker (Docker Desktop or Docker Engine + Compose plugin)
- Go `1.25+`
- .NET SDK `8.0` (needed for `product-service`, currently WIP)

## Quick Start (Full Stack via Docker)

From the repository root:

```bash
cd infra
docker compose up -d
```

Useful commands:

```bash
docker compose ps
docker compose logs -f
docker compose down
docker compose down -v
```

## Local URLs and Default Credentials

| Component | URL | Credentials |
| --- | --- | --- |
| Nginx gateway | http://localhost | - |
| User service | http://localhost:8081 | - |
| Product service (WIP) | http://localhost:8082 | - |
| Search service | http://localhost:8086 | - |
| Keycloak | http://localhost:8180 | `admin / admin123` |
| Grafana | http://localhost:3000 | `admin / grafana123` |
| Prometheus | http://localhost:9090 | - |
| Jaeger | http://localhost:16686 | - |
| pgAdmin | http://localhost:5050 | `admin@ecommerce.com / pgadmin123` |
| Kibana | http://localhost:5601 | - |
| Kafdrop | http://localhost:9000 | - |
| Elasticsearch | http://localhost:9200 | - |
| PostgreSQL | `localhost:5432` | `ecommerce / ecommerce123` |
| Redis | `localhost:6379` | password: `redis123` |
| Kafka | `localhost:9092` | - |

## Run Services Locally (Without Docker Build)

Run shared infra first (`cd infra && docker compose up -d`), then start services from separate terminals.

### user-service

```bash
cd services/user-service
cp .env.example .env
go run cmd/main.go
```

Default port: `8081`

### search-service

```bash
cd services/search-service
cp .env.example .env
go run cmd/main.go
```

Default port: `8086`

## API Surface (Current)

### Gateway-routed endpoints (`http://localhost`)

- `GET /health` (Nginx health)
- `GET /realms/*` (Keycloak passthrough)
- `GET /api/v1/products`
- `GET /api/v1/products/{id}`
- `GET /api/v1/categories`
- `POST|PUT|DELETE /api/v1/products...` (auth required)
- `/api/v1/users/*` (auth required)

Note: search is currently **not** routed via Nginx in `infra/nginx/conf.d/ecommerce.conf`.
Note: product routes are available in gateway config but the `.NET` service itself is still WIP.

### Direct service endpoints

- user-service (`:8081`)
- `GET /health`
- `GET /metrics`
- `GET /internal/auth/validate` (used internally by Nginx `auth_request`)
- `GET /api/v1/users/me`
- `PUT /api/v1/users/me`
- `GET|POST|PUT|DELETE|PATCH /api/v1/users/me/addresses...`

- search-service (`:8086`)
- `GET /search?q=<query>`

## Authentication

Protected routes are validated through:
1. Client sends bearer token to Nginx.
2. Nginx calls `user-service` endpoint `/internal/auth/validate`.
3. `user-service` introspects token against Keycloak.
4. Nginx forwards request with user headers to downstream service.

Token endpoint (via gateway):

```bash
curl -X POST http://localhost/realms/ecommerce/protocol/openid-connect/token \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=password' \
  -d 'client_id=ecommerce-app' \
  -d 'username=<username>' \
  -d 'password=<password>'
```

## Database and Migrations

- Databases are created by `infra/init-scripts/postgres/01-init-databases.sql`.
- `user-service` SQL migration files exist in `services/user-service/migrations` and are not auto-applied by the service itself.

## Tests

Executed and passing in this environment:

```bash
cd services/user-service
env GOCACHE=/tmp/ecommerce-go-build-cache GOMODCACHE=/tmp/ecommerce-go-mod-cache go test ./...

cd ../search-service
env GOCACHE=/tmp/ecommerce-go-build-cache GOMODCACHE=/tmp/ecommerce-go-mod-cache go test ./...
```

## Infrastructure Notes

- Compose file: `infra/docker-compose.yml`
- Keycloak realm import: `infra/keycloak/realm-export.json`
- K8s manifests: `infra/k8s`
- Terraform config: `infra/terraform`

## Known Gaps

- `order-service`, `payment-service`, and `notification-service` are placeholders only.
- `product-service` is currently WIP and may change without compatibility guarantees.
- `docs/openapi` does not yet contain service specs.
- Prometheus config includes scrape jobs for not-yet-implemented services.
