.PHONY: up down stop restart log

up:
	docker compose up --build --wait -d

down:
	docker compose down

stop:
	docker compose stop

restart:
	docker compose restart
	docker compose up --wait -d

log:
	docker compose logs -f
