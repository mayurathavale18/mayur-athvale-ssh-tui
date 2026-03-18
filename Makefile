BINARY_NAME=ssh-portfolio
BUILD_DIR=build
CMD=./cmd/server
DOCKER_IMAGE=ssh-portfolio

.PHONY: build run clean test lint docker-build docker-run docker-up docker-down build-all

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD)

run:
	go run $(CMD)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./... -v

lint:
	go vet ./...

# Docker
docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run --rm -p 2222:22 -v ssh-keys:/app/.ssh -v analytics:/app/data $(DOCKER_IMAGE)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Cross-compile
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD)

build-all: build-linux-amd64 build-linux-arm64
