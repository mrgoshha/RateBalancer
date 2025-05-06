SHELL := /bin/bash

COMPOSE_FILE := docker-compose.yaml
MIGRATE_CMD  := migrate -path ./migrations -database

.PHONY: up down migrate test test-db test-migrate test-clean

up:
	docker compose -f $(COMPOSE_FILE) up -d

down:
	docker compose -f $(COMPOSE_FILE) down

migrate:
	$(MIGRATE_CMD) "postgres://pguser:pgpwd@localhost:5432/rateLimiter?sslmode=disable" up

test:
	go test ./...
	$(MAKE) test-clean

test-db:
	docker run -d \
	  --name postgres-rate-limiter-test \
	  -e POSTGRES_USER=test \
	  -e POSTGRES_PASSWORD=test\
	  -e POSTGRES_DB=rateLimiterTest \
	  -p 5432:5432 \
	  postgres:latest

test-migrate:
	$(MIGRATE_CMD) "postgres://test:test@localhost:5432/rateLimiterTest?sslmode=disable" up

test-clean:
	docker rm -f postgres-rate-limiter-test

