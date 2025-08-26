
.PHONY: all backend frontend up down logs

all: backend frontend

backend:
	$(MAKE) -C backend all

frontend:
	docker compose -f frontend/docker-compose-frontend.yml up --build -d

up-all: backend frontend

down-all:
	docker compose -f backend/docker-compose-backend.yml down
	docker compose -f backend/docker-compose-databases.yml down
	docker compose -f backend/docker-compose-kafka.yml down
	docker compose -f frontend/docker-compose-frontend.yml down

logs-all:
	docker compose -f backend/docker-compose-backend.yml logs -f
	docker compose -f frontend/docker-compose-frontend.yml logs -f
