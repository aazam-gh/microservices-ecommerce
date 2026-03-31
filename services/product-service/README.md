# Product Service (.NET)

ASP.NET Core catalog service for the e-commerce platform.

## Local Run

1. Copy env vars: `cp .env.example .env`
2. Export env vars from `.env` (or run via Docker Compose)
3. Start service:
   - `dotnet run --project Ecommerce.ProductService.csproj`

Service port defaults to `8082`.

## Endpoints

- `GET /health`
- `GET /metrics`
- `GET /api/v1/products`
- `GET /api/v1/products/{id}`
- `POST /api/v1/products`
- `PUT /api/v1/products/{id}`
- `DELETE /api/v1/products/{id}`
- `GET /api/v1/categories`

`POST`, `PUT`, and `DELETE` routes are protected when called through Nginx.
