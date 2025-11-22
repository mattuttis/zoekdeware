.PHONY: help dev dev-down build build-all test lint proto clean k8s-dev k8s-prod k8s-down

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
