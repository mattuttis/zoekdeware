.PHONY: help dev dev-down build build-all test lint proto clean

help:
	@echo "Dating App - Available commands:"
	@echo ""
	@echo "  make dev          - Start development infrastructure (postgres, redis, nats)"
	@echo "  make dev-down     - Stop development infrastructure"
	@echo "  make build        - Build all services"
	@echo "  make test         - Run all tests"
	@echo "  make lint         - Run linter on all services"
	@echo "  make proto        - Generate protobuf code"
	@echo "  make clean        - Clean build artifacts"

dev:
	./scripts/dev.sh up

dev-down:
	./scripts/dev.sh down

build:
	./scripts/build-all.sh

build-%:
	./scripts/build-service.sh $*

test:
	cd backend && go test ./...

lint:
	cd backend/gateway && golangci-lint run
	cd backend/services/member && golangci-lint run
	cd backend/services/matching && golangci-lint run
	cd backend/services/messaging && golangci-lint run
	cd backend/services/notification && golangci-lint run
	cd backend/services/media && golangci-lint run
	cd backend/services/location && golangci-lint run

proto:
	./scripts/generate-proto.sh

clean:
	rm -rf backend/services/*/bin
	rm -rf backend/gateway/bin
	docker rmi $$(docker images -q dating-*) 2>/dev/null || true
