-- ============================================================
-- E-Commerce Platform — PostgreSQL Init Script
-- Creates a separate database for each microservice
-- File path: init-scripts/postgres/01-init-databases.sql
-- ============================================================

-- user-service database
CREATE DATABASE users_db;
GRANT ALL PRIVILEGES ON DATABASE users_db TO ecommerce;

-- product-service database
CREATE DATABASE products_db;
GRANT ALL PRIVILEGES ON DATABASE products_db TO ecommerce;

-- order-service database
CREATE DATABASE orders_db;
GRANT ALL PRIVILEGES ON DATABASE orders_db TO ecommerce;

-- payment-service database
CREATE DATABASE payments_db;
GRANT ALL PRIVILEGES ON DATABASE payments_db TO ecommerce;

-- notification-service database
CREATE DATABASE notifications_db;
GRANT ALL PRIVILEGES ON DATABASE notifications_db TO ecommerce;

-- Keycloak database
CREATE DATABASE keycloak_db;
GRANT ALL PRIVILEGES ON DATABASE keycloak_db TO ecommerce;
