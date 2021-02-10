BIN := "./bin"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -mod=mod -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -mod=mod -v -o $(BIN) -ldflags "$(LDFLAGS)" ./calendar_scheduler
	go build -mod=mod -v -o $(BIN) -ldflags "$(LDFLAGS)" ./calendar_sender

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f calendar_scheduler/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.36.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

generate:
	mkdir -p internal/server/pb
	protoc --proto_path=api/ --go_out=internal/server/pb  --go-grpc_out=internal/server/pb  api/*.proto

lint:
	golangci-lint run ./...

lint-fix:
	gofmt -w ./..
	gci -w ./..

up:
	sudo docker-compose up -d --build

down:
	docker-compose down

integration-tests:
