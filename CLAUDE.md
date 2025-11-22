# ZoekDeware - Project Context

## Project Overview
Dating application with:
- **Backend**: Go 1.24 microservices (gateway, member service) with DDD/Event Sourcing
- **Infrastructure**: Kubernetes with dev/prod overlays, PostgreSQL, NATS
- **Mobile**: iOS app (SwiftUI, iOS 16+ for iPhone 14+)

## Tech Stack
- Go 1.24 (matches golangci-lint v1.64.8)
- Docker with distroless images
- Kubernetes with Kustomize overlays
- gRPC for inter-service communication
- REST API via gateway
- SwiftUI for iOS

## Key Commands
```bash
make dev          # Start Kubernetes dev environment
make dev-down     # Stop dev environment
make build        # Build all Docker images
make lint         # Run golangci-lint
make test         # Run tests
```

## Architecture
```
backend/
├── gateway/           # REST API gateway (port 8080)
├── services/
│   └── member/        # Member service (gRPC port 9090)
└── shared/            # Shared packages

mobile/ios/ZoekDeware/ # iOS app
infrastructure/kubernetes/ # K8s manifests with overlays
```

## Current State (Nov 2025)
- CI pipeline: Go 1.24 + golangci-lint v1.64.8
- Docker builds working
- Kubernetes dev environment operational
- iOS app connects to gateway
- Service handlers return 501 (not yet implemented)

## Pending Work
- Implement member service handlers
- Add remaining services: matching, messaging, notification, media, location
- Complete iOS app features
- Add PostgreSQL and NATS integration

## Notes
- Dockerfiles use Go 1.25-alpine (backward compatible with go.mod 1.24)
- golangci-lint must match Go version (v1.64.8 built with Go 1.24)
- Coverage flags removed from CI (covdata tool issues)
