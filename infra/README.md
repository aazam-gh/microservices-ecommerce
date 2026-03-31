# E-Commerce Platform — Infrastructure

## Folder Structure

```
infra/
├── docker-compose.yml
├── init-scripts/
│   └── postgres/
│       └── 01-init-databases.sql     # Automatically generates databases
├── nginx/
│   ├── nginx.conf                    # Main Nginx config
│   └── conf.d/
│       └── ecommerce.conf            # Routing + auth_request
├── prometheus/
│   └── prometheus.yml                # Scrape config
├── grafana/
│   ├── provisioning/                 # Automatic datasource binding
│   └── dashboards/                   # Dashboard JSONs
├── pgadmin/
│   └── servers.json                  # Automatic pgAdmin connection
└── keycloak/
    └── realm-export.json             # Realm automatically imported
```

## Startup

```bash
# Bring up the entire infrastructure
docker-compose up -d

# Follow the logs
docker-compose logs -f

# Logs for a specific service
docker-compose logs -f kafka

# Stop
docker-compose down

# Stop and remove volumes (clean start)
docker-compose down -v
```

## Service Addresses

| Service          | URL                          | Credentials                   |
|------------------|------------------------------|-------------------------------|
| Nginx            | http://localhost:80          | —                             |
| Product Service  | http://localhost:8082        | —                             |
| Keycloak Admin   | http://localhost:8180        | admin / admin123              |
| Kafdrop          | http://localhost:9000        | —                             |
| Kibana           | http://localhost:5601        | —                             |
| Grafana          | http://localhost:3000        | admin / grafana123            |
| Jaeger UI        | http://localhost:16686       | —                             |
| pgAdmin          | http://localhost:5050        | admin@ecommerce.com / pgadmin123 |
| PostgreSQL       | localhost:5432               | ecommerce / ecommerce123      |
| Redis            | localhost:6379               | redis123                      |
| Kafka            | localhost:9092               | —                             |
| Elasticsearch    | http://localhost:9200        | —                             |

## Databases

Databases automatically created inside PostgreSQL:

| Database         | Service               |
|------------------|-----------------------|
| users_db         | user-service          |
| products_db      | product-service       |
| orders_db        | order-service         |
| payments_db      | payment-service       |
| notifications_db | notification-service  |
| keycloak_db      | Keycloak              |

## Kafka Topics

Automatically created topics:

- product.created / product.updated / product.deleted
- order.created / order.confirmed / order.cancelled
- stock.reserve / stock.reserved / stock.reserve.failed / stock.release
- payment.process / payment.completed / payment.failed
- notification.send

## After Keycloak Setup

1. Go to http://localhost:8180/admin (admin / admin123)
2. The `ecommerce` realm should already be imported
3. If not: Realm → Create → import `keycloak/realm-export.json`
4. Use the `ecommerce-app` client for the frontend (public)
5. Use the `ecommerce-service` client for backend services (requires secret)

## Nginx + Keycloak Auth Flow

```
Client
  → Nginx :80
  → auth_request → user-service:8081/internal/auth/validate
      → Keycloak:8180/realms/ecommerce/protocol/openid-connect/token/introspect
  → Downstream service (with X-User-ID, X-User-Role headers)
```

## Adding Your Services

After writing your services, add them to docker-compose.yml:

```yaml
user-service:
  build: ../services/user-service
  container_name: ecommerce-user-service
  ports:
    - "8081:8081"
  environment:
    DB_HOST: postgres
    DB_PORT: 5432
    DB_NAME: users_db
    DB_USER: ecommerce
    DB_PASSWORD: ecommerce123
    REDIS_ADDR: redis:6379
    REDIS_PASSWORD: redis123
    KEYCLOAK_URL: http://keycloak:8180
    KEYCLOAK_REALM: ecommerce
    JAEGER_ENDPOINT: http://jaeger:4318/v1/traces
  depends_on:
    postgres:
      condition: service_healthy
    redis:
      condition: service_healthy
    keycloak:
      condition: service_healthy
  networks:
    - ecommerce-net
```

`product-service` already ships as an ASP.NET Core service in `services/product-service` and is wired in `docker-compose.yml` with `products_db`.

search-service:
  build:
    context: ../services/search-service
    dockerfile: Dockerfile
  container_name: ecommerce-search-service
  restart: unless-stopped
  environment:
    APP_PORT: 8086
    ELASTICSEARCH_URL: http://elasticsearch:9200
    ELASTICSEARCH_INDEX: products
    ELASTICSEARCH_TIMEOUT: 5s
    REDIS_ENABLED: "true"
    REDIS_ADDR: redis:6379
    REDIS_PASSWORD: redis123
    REDIS_DB: 0
    REDIS_KEY_PREFIX: search-service
    REDIS_CACHE_TTL: 2m
  ports:
    - "8086:8086"
  depends_on:
    elasticsearch:
      condition: service_healthy
    redis:
      condition: service_healthy
  networks:
    - ecommerce-net
