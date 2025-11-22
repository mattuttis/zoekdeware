#!/bin/bash
set -e

TAG=${1:-latest}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

SERVICES="gateway member"

for service in $SERVICES; do
    echo "========================================"
    echo "Building $service..."
    echo "========================================"
    "$SCRIPT_DIR/build-service.sh" "$service" "$TAG"
done

echo ""
echo "All services built successfully!"
docker images | grep "dating-"
