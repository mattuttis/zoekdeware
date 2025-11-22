# Dating App

A microservices-based dating application backend built with Go, following DDD and Event Sourcing patterns.

## Architecture

```
┌─────────────┐     ┌─────────────────────────────────────────────┐
│   Mobile    │────▶│              API Gateway                    │
│   Clients   │     │         (REST/WebSocket)                    │
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
│   │   ├── cmd/server/
│   │   └── internal/
│   │       ├── config/
│   │       ├── handlers/
│   │       ├── middleware/
│   │       └── router/
│   ├── services/
│   │   └── member/           # Member service (DDD + Event Sourcing)
│   │       ├── cmd/server/
│   │       └── internal/
│   │           ├── domain/
│   │           │   ├── aggregate/
│   │           │   ├── commands/
│   │           │   ├── events/
│   │           │   ├── repository/
│   │           │   └── valueobject/
│   │           ├── application/
│   │           └── infrastructure/
│   │               └── eventstore/
│   ├── shared/               # Shared packages
│   │   └── pkg/
│   │       ├── auth/
│   │       ├── config/
│   │       ├── errors/
│   │       ├── eventstore/
│   │       └── messaging/
│   └── go.work               # Go workspace
├── infrastructure/
│   ├── docker-compose.yml
│   └── kubernetes/
│       ├── base/
│       └── overlays/
│           ├── dev/
│           └── prod/
├── scripts/
│   ├── build-service.sh
│   ├── build-all.sh
│   ├── dev.sh
│   └── generate-proto.sh
└── .github/workflows/
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Make

### Development

1. Start infrastructure (PostgreSQL, Redis, NATS):
```bash
make dev
```

2. Run the member service:
```bash
cd backend/services/member
go run ./cmd/server
```

3. Run the gateway:
```bash
cd backend/gateway
go run ./cmd/server
```

### Building

Build a specific service:
```bash
make build-member
make build-gateway
```

Build all services:
```bash
make build
```

### Running with Docker Compose

```bash
cd infrastructure
docker-compose up
```

## API Documentation

- REST API: See `api/openapi/gateway.yaml`
- gRPC: See `api/proto/member/v1/member.proto`

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
