package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	memberv1 "github.com/mattuttis/inetcontrol/zoekdeware/api/proto/member/v1"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/config"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/handlers"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/router"
)

func main() {
	cfg := config.Load()

	// Connect to member service
	memberConn, err := grpc.NewClient(
		cfg.MemberServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to member service: %v", err)
	}
	defer memberConn.Close()

	memberClient := memberv1.NewMemberServiceClient(memberConn)

	// Create handlers with gRPC clients
	h := handlers.NewHandlers(memberClient, cfg.JWTSecret)

	// Create router
	r := router.New(cfg, h)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API Gateway listening on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
