# ZoekDeware

A microservices-based dating application built with Go, following DDD and Event Sourcing patterns.

## Architecture

```
┌─────────────┐     ┌─────────────────────────────────────────────┐
│   iOS App   │────▶│              API Gateway                    │
│  (SwiftUI)  │     │         (REST/WebSocket)                    │
└─────────────┘     └─────────────────────────────────────────────┘
                                     │
                                     ▼
                              ┌──────────┐
                              │  Member  │
                              │ Service  │
                              └──────────┘
                                     │
                             ┌───────┴───────┐
                             │  Event Bus    │
                             │   (NATS)      │
                             └───────────────┘
```

## Services

| Service | Port (HTTP) | Port (gRPC) | Status | Description |
|---------|-------------|-------------|--------|-------------|
| Gateway | 8000 | - | Active | API Gateway for mobile clients |
| Member | 8080 | 9090 | Active | User registration, profiles |
| Matching | - | - | Planned | Swipes, likes, matches |
| Messaging | - | - | Planned | Chat conversations |
| Notification | - | - | Planned | Push notifications |
| Media | - | - | Planned | Photo uploads |
| Location | - | - | Planned | Geolocation |

## Project Structure

```
.
├── api/
│   ├── proto/
│   │   └── member/v1/        # Member service gRPC definitions
│   └── openapi/
│       └── gateway.yaml      # REST API spec
├── backend/
│   ├── gateway/              # API Gateway
│   ├── services/
│   │   └── member/           # Member service (DDD + Event Sourcing)
│   ├── shared/               # Shared packages
│   └── go.work               # Go workspace
├── mobile/
│   └── ios/                  # iOS app (SwiftUI, iOS 16+)
│       ├── ZoekDeware/
│       ├── project.yml       # XcodeGen configuration
│       └── Makefile
├── infrastructure/
│   ├── docker-compose.yml
│   └── kubernetes/
│       ├── base/             # Base Kubernetes manifests
│       └── overlays/
│           ├── dev/          # Development environment
│           └── prod/         # Production environment
├── scripts/
└── .github/workflows/
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Kubernetes (Docker Desktop or minikube)
- Make

### Option 1: Docker Compose (Simple)

```bash
# Start PostgreSQL, Redis, NATS
make dev

# Run services locally
cd backend/services/member && go run ./cmd/server
cd backend/gateway && go run ./cmd/server
```

### Option 2: Kubernetes (Recommended)

```bash
# Build Docker images
make build

# Deploy to local Kubernetes
make k8s-dev

# Access the API
curl http://localhost/health
```

**Kubernetes Dev Environment:**
- Gateway: http://localhost (LoadBalancer)
- PostgreSQL: dev-postgres:5432
- NATS: dev-nats:4222

Check pod status:
```bash
kubectl get pods -n dating-app-dev
```

### iOS App

Requires Xcode 15+ and XcodeGen:

```bash
cd mobile/ios
make setup   # Install XcodeGen and generate project
make open    # Open in Xcode
```

See [mobile/ios/README.md](mobile/ios/README.md) for details.

## Building

```bash
# Build all Docker images
make build

# Build specific service
make build-member
make build-gateway
```

## Kubernetes Environments

| Environment | Namespace | Replicas | Resources |
|-------------|-----------|----------|-----------|
| Dev | dating-app-dev | 1 | 200m CPU, 256Mi |
| Prod | dating-app | 2+ | 500m CPU, 512Mi |

Deploy:
```bash
make k8s-dev   # Development
make k8s-prod  # Production
```

## API Documentation

- REST API: `api/openapi/gateway.yaml`
- gRPC: `api/proto/member/v1/member.proto`

## Domain Model (Member Service)

The member service follows DDD patterns:

- **Aggregate**: `Member` - manages user identity and profile
- **Value Objects**: `Email`, `Profile`, `Gender`
- **Domain Events**: `MemberRegistered`, `ProfileUpdated`, `MemberActivated`
- **Commands**: `RegisterMember`, `UpdateProfile`, `ActivateMember`

## Testing

```bash
make test
```

## License

Proprietary
