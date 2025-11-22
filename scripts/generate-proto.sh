#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
API_DIR="$ROOT_DIR/api/proto"

if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Install with: brew install protobuf"
    exit 1
fi

if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

SERVICES="member matching messaging notification media location"

for service in $SERVICES; do
    PROTO_DIR="$API_DIR/$service/v1"
    if [ -d "$PROTO_DIR" ]; then
        echo "Generating Go code for $service..."
        protoc --proto_path="$API_DIR" \
            --go_out="$ROOT_DIR" --go_opt=paths=source_relative \
            --go-grpc_out="$ROOT_DIR" --go-grpc_opt=paths=source_relative \
            "$PROTO_DIR"/*.proto
    fi
done

echo "Proto generation complete!"
