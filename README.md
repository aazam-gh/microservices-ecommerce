# 📦 E-Commerce Platform

Distributed, saga-driven e-commerce backend built with Go, .NET, Kafka, and centralized Keycloak authentication. Each service owns its database, follows event-driven best practices (Outbox + Saga choreography), and exposes a clear HTTP contract behind Nginx.

## Highlights
1. Microservices map: user, catalog, orders, payments, search, and notifications wired through Kafka topics and idempotent handlers.
2. Resilient patterns: transactional outbox tables, saga choreography for distributed consistency, retry-friendly compensation, and circuit breakers at the edge.
3. Observability-first: structured Zap logs, Prometheus metrics, distributed traces in Jaeger, and centralized dashboards.

## Architecture & Stack
| Layer | Technology |
| --- | --- |
| Language | Go 1.25+ and .NET 8 (product-service) |
| HTTP | Fiber v2 (Go services), ASP.NET Core (product-service) |
| Persistence | PostgreSQL 16 (one DB per service) |
| Search | Elasticsearch 8 |
| Messaging | Apache Kafka |
| Cache | Redis 7 |
| Auth | Keycloak 24 (OAuth2 / OIDC) |
| Edge | Nginx + auth_request |
| Packaging | Docker Compose for local, Kubernetes for production |
| CI/CD | GitHub Actions + Terraform |
| Observability | Prometheus, Grafana, Jaeger |

## Quick Start (Local)
1. Install Go 1.25+, .NET 8 SDK, Docker Desktop, and Git.
2. Clone the repo and enter it.
3. Launch the shared infrastructure from `infra/docker-compose.yml` and wait ~60 seconds for Keycloak.
4. Confirm services with `docker-compose ps` and hit http://localhost:8180/realms for Keycloak readiness.
5. Tooling credentials: Keycloak admin/admin123, Grafana admin/grafana123, pgAdmin admin@ecommerce.com/pgadmin123.

## Running a Service Locally
1. `cd services/<service>` (for example `services/user-service`, `services/product-service`, or `services/search-service`).
2. `cp .env.example .env` and adapt secrets if needed.
3. Start the service:
   - Go services: `go run cmd/main.go`
   - product-service (.NET): `dotnet run --project Ecommerce.ProductService.csproj`
4. Use a Keycloak token (`/realms/ecommerce/protocol/openid-connect/token`) for protected endpoints, passing the bearer token on requests.
5. `user-service` now uses Redis for Keycloak token introspection caching; tune `REDIS_AUTH_CACHE_TTL`, `REDIS_KEY_PREFIX`, or disable it with `REDIS_ENABLED=false`.
6. `search-service` proxies Elasticsearch via `/search?q=…` and caches successful responses in Redis (`REDIS_CACHE_TTL`, `REDIS_KEY_PREFIX`, `ELASTICSEARCH_INDEX`).

## Service Map
| Service | Port | Focus |
| --- | --- | --- |
| user-service | 8081 | Profiles, addresses, Keycloak token validation |
| product-service | 8082 | Catalog API on ASP.NET Core + PostgreSQL |
| order-service | 8083 | Order lifecycle, saga orchestration, compensations |
| payment-service | 8084 | Payment provider integrations, retries, refunds |
| notification-service | 8085 | Async email/SMS via Kafka consumers |
| search-service | 8086 | Elasticsearch queries (cached in Redis) and result delivery |
| Nginx | 80/443 | JWT validation, request routing, fallback pages |
| Keycloak | 8180 | Identity provider, realm configuration, user mgmt |

## Infrastructure Notes
- Docker Compose in `infra/` wires PostgreSQL instances, Kafka (with Schema Registry), Redis, Keycloak, Grafana, Prometheus, Jaeger, and observability tooling for local development.
- `product-service` is implemented in .NET 8 and bootstraps its schema on startup (`EnsureCreated`) against `products_db`.
- Docker Compose now also starts `search-service`, which depends on Elasticsearch + Redis and caches query results for faster autocomplete/lookup responses.
- Kubernetes manifests (`infra/k8s/`) describe deployments, services, secrets, and ingress rules for production, with Terraform managing cloud networking, storage, and clusters.
- Keycloak realm exports live in `infra/keycloak/`; run `infra/keycloak/import.sh` (or the equivalent script) after an upgrade to re-import the realm.
- Terraform modules live under `infra/terraform/` and cover cloud resources for databases, Kafka, and load balancing.

## Observability & Monitoring
- Logs stream structured JSON via `go.uber.org/zap` with fields `service_name`, `request_id`, and `user_id` and are aggregated by Loki/ELK in production.
- Prometheus scrapes `/metrics` on every service. Dashboards in Grafana cover latency, error rates, DB pools, and Kafka lag.
- Jaeger collects OpenTelemetry traces via OTLP HTTP (`http://jaeger:4318/v1/traces`) to visualize distributed requests.

## API & Docs
- OpenAPI specs live under `docs/openapi/` for each service; use Swagger or any spec-aware client to validate contracts.
- Public endpoints (products, search, categories) accept anonymous requests. Private endpoints require Keycloak-issued JWTs tied to `ecommerce-app` clients.
- product-service endpoints:
  - `GET /api/v1/products` and `GET /api/v1/products/{id}` (public)
  - `GET /api/v1/categories` (public)
  - `POST/PUT/DELETE /api/v1/products...` (protected by Nginx auth_request)
- Authentication flow: POST to `/realms/ecommerce/protocol/openid-connect/token` with client credentials, then include `Authorization: Bearer <token>` in API calls.

## Development Workflow
- CI (GitHub Actions) path filters trigger only the touched service pipeline. Go services run `go vet`, `golangci-lint`, and `go test ./... -race -cover`; product-service runs `dotnet build`/`dotnet test`; all services build Docker images and push on `main`.
- Run local smoke tests by hitting `/api/v1/products`, `/api/v1/orders`, and `/api/v1/users/me` with the appropriate scopes.
- Kafka topics are under `infra/kafka/topics.yml`; use `kafdrop` (http://localhost:9000) to inspect message flows.

## Contributing & Testing
- Add unit tests near the behavior you change. Use mocks for repositories and Kafka producers.
- Integration tests rely on the infra stack started via Docker Compose; run them inside the relevant service directory (`go test ./...` for Go services, `dotnet test` for product-service).
- For schema changes, update SQL migrations in `services/<service>/migrations` and ensure they are idempotent.

## License
MIT License. See `LICENSE`.
