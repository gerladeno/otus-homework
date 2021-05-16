BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
PG_USER := "calendar"
PG_PASSWORD := "calendar"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
DOCKER_GATEWAY_HOST := $(shell hostname -I | awk '{print $$1}')

generate:
	go generate internal/server/grpc/server.go

build_calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build_scheduler:
	go build -v -o ./bin/scheduler ./cmd/scheduler

build_sender:
	go build -v -o ./bin/sender ./cmd/sender

start_db:
	PG_USER=$(PG_USER) PG_PASSWORD=$(PG_PASSWORD) docker-compose -f deployments/docker-compose.yml up calendar-postgres

start_rabbit:
	PG_USER=$(PG_USER) PG_PASSWORD=$(PG_PASSWORD) docker-compose -f deployments/docker-compose.yml up calendar-rabbit

run_calendar: build_calendar
	$(BIN) --config ./configs/calendar_config.json

run_scheduler: build_scheduler
	./cmd/scheduler --config ./configs/scheduler_config.json

run_sender: build_sender
	./cmd/sender --config ./configs/sender_config.json

build_images: generate
	docker-compose -f deployments/docker-compose.yml build

up: build_images
	DOCKER_GATEWAY_HOST=$(DOCKER_GATEWAY_HOST) PG_USER=$(PG_USER) PG_PASSWORD=$(PG_PASSWORD) \
	docker-compose -f deployments/docker-compose.yml up -d \
 	calendar-postgres calendar-rabbit calendar-calendar calendar-scheduler calendar-sender

down:
	docker-compose -f deployments/docker-compose.yml down

version: build_calendar
	$(BIN) version

test:
	go test -race -count=50 ./internal/server/grpc ./internal/server/http ./internal/storage/memory ./internal/storage/sql

integration-tests: up
	docker-compose -f deployments/docker-compose.yml run --service-ports calendar-integration-tests
	docker-compose -f deployments/docker-compose.yml down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
