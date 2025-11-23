.PHONY: help dev dev-down build build-all test lint proto clean k8s-dev k8s-prod k8s-down \
	migrate-up migrate-down migrate-create migrate-version migrate-force

help:
	@echo "ZoekDeware - Available commands:"
	@echo ""
	@echo "  make dev          - Start development infrastructure (docker-compose)"
	@echo "  make dev-down     - Stop development infrastructure"
	@echo "  make build        - Build all Docker images"
	@echo "  make test         - Run all tests"
	@echo "  make lint         - Run linter on all services"
	@echo "  make proto        - Generate protobuf code"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make k8s-dev      - Deploy to Kubernetes (dev)"
	@echo "  make k8s-prod     - Deploy to Kubernetes (prod)"
	@echo "  make k8s-down     - Remove Kubernetes deployment (dev)"
	@echo ""
	@echo "Migrations:"
	@echo "  make migrate-up SERVICE=member         - Run migrations up"
	@echo "  make migrate-down SERVICE=member       - Rollback last migration"
	@echo "  make migrate-create SERVICE=member NAME=xyz - Create new migration"
	@echo "  make migrate-version SERVICE=member    - Show current version"

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
	cd backend/gateway && ~/go/bin/golangci-lint run
	cd backend/services/member && ~/go/bin/golangci-lint run

proto:
	./scripts/generate-proto.sh

clean:
	rm -rf backend/services/*/bin
	rm -rf backend/gateway/bin
	docker rmi $$(docker images -q dating-*) 2>/dev/null || true

# Kubernetes targets
k8s-dev:
	kubectl apply -k infrastructure/kubernetes/overlays/dev

k8s-prod:
	kubectl apply -k infrastructure/kubernetes/overlays/prod

k8s-down:
	kubectl delete -k infrastructure/kubernetes/overlays/dev --ignore-not-found

# Database migrations
# Usage: make migrate-up SERVICE=member DATABASE_URL=postgres://user:pass@localhost:5432/db?sslmode=disable
MIGRATIONS_PATH = backend/services/$(SERVICE)/migrations

migrate-up:
	@if [ -z "$(SERVICE)" ]; then echo "SERVICE is required (e.g., SERVICE=member)"; exit 1; fi
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	@if [ -z "$(SERVICE)" ]; then echo "SERVICE is required (e.g., SERVICE=member)"; exit 1; fi
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-create:
	@if [ -z "$(SERVICE)" ]; then echo "SERVICE is required (e.g., SERVICE=member)"; exit 1; fi
	@if [ -z "$(NAME)" ]; then echo "NAME is required (e.g., NAME=create_users)"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

migrate-version:
	@if [ -z "$(SERVICE)" ]; then echo "SERVICE is required (e.g., SERVICE=member)"; exit 1; fi
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

migrate-force:
	@if [ -z "$(SERVICE)" ]; then echo "SERVICE is required (e.g., SERVICE=member)"; exit 1; fi
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(VERSION)
