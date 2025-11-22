#!/bin/bash
set -e

SERVICE=$1
TAG=${2:-latest}

if [ -z "$SERVICE" ]; then
    echo "Usage: $0 <service-name> [tag]"
    echo "Available services: gateway, member, matching, messaging, notification, media, location"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

if [ "$SERVICE" = "gateway" ]; then
    DOCKERFILE="backend/gateway/Dockerfile"
else
    DOCKERFILE="backend/services/$SERVICE/Dockerfile"
fi

if [ ! -f "$ROOT_DIR/$DOCKERFILE" ]; then
    echo "Error: Dockerfile not found for service '$SERVICE'"
    exit 1
fi

echo "Building dating-$SERVICE:$TAG..."
docker build -f "$ROOT_DIR/$DOCKERFILE" -t "dating-$SERVICE:$TAG" "$ROOT_DIR"

echo "Successfully built dating-$SERVICE:$TAG"
