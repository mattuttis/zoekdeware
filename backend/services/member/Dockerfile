FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.work and all modules it references
COPY backend/go.work backend/go.work
COPY backend/shared/ backend/shared/
COPY backend/gateway/ backend/gateway/
COPY backend/services/member/ backend/services/member/

WORKDIR /app/backend/services/member

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12

COPY --from=builder /server /server

EXPOSE 8080 9090

CMD ["/server"]
