BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate internal/server/grpc/server.go

build_calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build_scheduler:
	go build -v -o ./bin/scheduler ./cmd/scheduler

build_sender:
	go build -v -o ./bin/sender ./cmd/sender

start_db:
	docker-compose -f deployments/docker-compose.yaml up postgres-calendar

start_rabbit:
	docker-compose -f deployments/docker-compose.yaml up rabbit-calendar

run_calendar: build_calendar
	$(BIN) --config ./configs/calendar_config.json

run_scheduler: build_scheduler
	./cmd/scheduler --config ./configs/scheduler_config.json

run_sender: build_sender
	./cmd/sender --config ./configs/sender_config.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build_calendar
	$(BIN) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
