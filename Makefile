SHELL := /bin/sh

-include .env

COMPOSE ?= docker compose
API_PORT ?= 18180
ADMIN_PORT ?= 18181
API_BASE_URL ?= http://localhost:$(API_PORT)
ADMIN_BASE_URL ?= http://localhost:$(ADMIN_PORT)
WAIT_RETRIES ?= 60
WAIT_SECONDS ?= 2

.PHONY: build up infra-up app-up migrate migrate-status ready wait-api wait-admin down stop restart reset log logs ps

build:
	$(COMPOSE) build api admin

infra-up:
	$(COMPOSE) up -d --wait postgres redis

migrate:
	$(COMPOSE) run --rm --no-deps api migrate up

migrate-status:
	$(COMPOSE) run --rm --no-deps api migrate status

app-up:
	$(COMPOSE) up -d --wait api admin

wait-api:
	@attempt=1; \
	while [ $$attempt -le $(WAIT_RETRIES) ]; do \
		if curl -fsS "$(API_BASE_URL)/api/v1/healthz" >/dev/null && \
			curl -fsS "$(API_BASE_URL)/api/v1/readyz" >/dev/null && \
			curl -fsS "$(API_BASE_URL)/api/v1/auth/bootstrap" >/dev/null; then \
			echo "API ready at $(API_BASE_URL)"; \
			exit 0; \
		fi; \
		echo "Waiting for API readiness ($$attempt/$(WAIT_RETRIES))..."; \
		attempt=`expr $$attempt + 1`; \
		sleep $(WAIT_SECONDS); \
	done; \
	echo "API readiness failed at $(API_BASE_URL)"; \
	exit 1

wait-admin:
	@attempt=1; \
	while [ $$attempt -le $(WAIT_RETRIES) ]; do \
		if curl -fsS "$(ADMIN_BASE_URL)/" >/dev/null; then \
			echo "Admin ready at $(ADMIN_BASE_URL)"; \
			exit 0; \
		fi; \
		echo "Waiting for Admin readiness ($$attempt/$(WAIT_RETRIES))..."; \
		attempt=`expr $$attempt + 1`; \
		sleep $(WAIT_SECONDS); \
	done; \
	echo "Admin readiness failed at $(ADMIN_BASE_URL)"; \
	exit 1

ready: wait-api wait-admin
	@echo "Finance Hub is browser-ready at $(ADMIN_BASE_URL)"

up: build infra-up migrate app-up ready

down:
	$(COMPOSE) down --remove-orphans

stop:
	$(COMPOSE) stop

restart: down up

reset:
	$(COMPOSE) down -v --remove-orphans

logs:
	$(COMPOSE) logs -f

log: logs

ps:
	$(COMPOSE) ps
