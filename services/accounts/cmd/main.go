package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/oarthurmonteiro/payment-system/services/accounts/gen/accounts/v1"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/database"
	handler "github.com/oarthurmonteiro/payment-system/services/accounts/internal/handler/grpc"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}
	
    if err := database.RunMigrations(dbURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
    }
	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

    pool, err := database.NewPostgresPool(ctx, dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer pool.Close()

    lis, err := net.Listen("tcp", ":50051") // Porta padrão gRPC
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	repo := database.NewPostgresRepository(pool)
	uc := usecase.NewOnboardingUseCase(repo)
	onboardingHandler := handler.NewOnboardingHandler(uc)

	// 4. Iniciar o gRPC Server
	server := grpc.NewServer()
	pb.RegisterOnboardingServiceServer(server, onboardingHandler)
	if os.Getenv("ENV") != "production" {
		reflection.Register(server)
	}
	
	// Rodar o servidor em uma Goroutine (para não travar o Graceful Shutdown)
	go func() {
		log.Printf("gRPC server listsening on %v", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Aguardar sinal de interrupção (Ctrl+C ou SIGTERM do Docker/K8s)
	<-ctx.Done()
	log.Println("Shutting down gRPC server gracefully...")
	server.GracefulStop()
	log.Println("Server stopped")
}