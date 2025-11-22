.PHONY: build run test lint proto docker

SERVICE_NAME := member
BINARY_NAME := server

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

proto:
	protoc --go_out=. --go-grpc_out=. api/proto/*.proto

docker:
	cd ../../.. && docker build -f backend/services/$(SERVICE_NAME)/Dockerfile -t dating-$(SERVICE_NAME) .

docker-run:
	docker run -p 8080:8080 -p 9090:9090 dating-$(SERVICE_NAME)

clean:
	rm -rf bin/ coverage.out coverage.html
