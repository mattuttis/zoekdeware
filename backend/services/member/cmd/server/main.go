package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	memberv1 "github.com/mattuttis/inetcontrol/zoekdeware/api/proto/member/v1"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/application"
	grpchandler "github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/interfaces/grpc"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/infrastructure/persistence"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Database connection
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/member?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("warning: database not available: %v", err)
	}

	// Initialize repository and service
	repo := persistence.NewPostgresMemberRepository(db)
	memberService := application.NewMemberService(repo, nil) // eventStore is optional for now

	// Initialize gRPC handler
	memberHandler := grpchandler.NewMemberHandler(memberService)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	memberv1.RegisterMemberServiceServer(grpcServer, memberHandler)
	reflection.Register(grpcServer) // Enable reflection for grpcurl

	// Health check HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	httpServer := &http.Server{
		Addr:         getEnv("HTTP_ADDR", ":8080"),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start gRPC server
	grpcAddr := getEnv("GRPC_ADDR", ":9090")
	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("gRPC server listening on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("HTTP server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcServer.GracefulStop()
	_ = httpServer.Shutdown(ctx)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
