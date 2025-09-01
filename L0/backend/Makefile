include .env
export

SHELL := /bin/bash

PROJECT_NAME := order_service
DATABASES := docker-compose-databases.yml
BACKEND := docker-compose-backend.yml
KAFKA := docker-compose-kafka.yml

COMPOSE_DATABASE := docker compose -f $(DATABASES)
COMPOSE_BACKEND := docker compose -f $(BACKEND)
COMPOSE_KAFKA := docker compose -f $(KAFKA)

DIRS = order_service

.PHONY: all
all: down-all network-all build-all up-all

.PHONY: down-all
down-all: down-backend down-databases down-kafka

.PHONY: build-all
build-all: build-databases build-kafka build-backend

.PHONY: up-all
up-all: up-databases up-kafka up-backend

# -------------------------------------------NETWORK-------------------------------------------

.PHONY: network-all
network-all: hard-rm-network create-network

.PHONY: create-network
create-network:
	@docker network create order_backend
	@docker network create order_frontend

.PHONY: rm-network
rm-network:
	@docker network rm order_backend || true
	@docker network rm order_frontend || true

.PHONY: hard-rm-network
hard-rm-network:
	@docker rm -f $$(docker ps -aq) || true
	@docker volume prune -f || true
	@docker network prune -f || true
	@docker network rm order_backend order_frontend || true

# -------------------------------------------KAFKA-------------------------------------------

.PHONY: up-kafka
up-kafka:
	$(COMPOSE_KAFKA) up -d --build
	@make kafka-topic-init

.PHONY: build-kafka
build-kafka:
	$(COMPOSE_KAFKA) build --no-cache --force-rm

.PHONY: down-kafka
down-kafka:
	$(COMPOSE_KAFKA) down || true

.PHONY: kafka-topic-init
kafka-topic-init:
	@docker exec kafka ./opt/kafka/bin/kafka-topics.sh --create --topic add-order --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1 --if-not-exists

# -------------------------------------------DATABASE-------------------------------------------

.PHONY: down-databases
down-databases:
	$(COMPOSE_DATABASE) down || true

.PHONY: build-databases
build-databases:
	$(COMPOSE_DATABASE) pull

.PHONY: up-databases
up-databases:
	$(COMPOSE_DATABASE) up -d --build

# -------------------------------------------SERVICES-------------------------------------------

.PHONY: down-backend
down-backend:
	$(COMPOSE_BACKEND) down || true

.PHONY: build-backend
build-backend:
	$(COMPOSE_BACKEND) build

.PHONY: up-backend
up-backend:
	$(COMPOSE_BACKEND) up --build -d

# -------------------------------------------FILES------------------------------------------- 

.PHONY: update-env
update-env:
	@for dir in $(DIRS); do \
	    rm -f $$dir/.env; \
	    cp .env $$dir/.env; \
	done
