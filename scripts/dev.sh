#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$ROOT_DIR/infrastructure"

case "${1:-up}" in
    up)
        echo "Starting development environment..."
        docker-compose up -d postgres redis nats
        echo "Waiting for services to be ready..."
        sleep 5
        echo "Infrastructure ready!"
        echo ""
        echo "PostgreSQL: localhost:5432"
        echo "Redis: localhost:6379"
        echo "NATS: localhost:4222"
        ;;
    down)
        echo "Stopping development environment..."
        docker-compose down
        ;;
    logs)
        docker-compose logs -f ${2:-}
        ;;
    ps)
        docker-compose ps
        ;;
    *)
        echo "Usage: $0 {up|down|logs|ps}"
        exit 1
        ;;
esac
