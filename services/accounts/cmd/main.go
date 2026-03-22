package main

import (
	"context"
	"log"
	"os"

	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/database"
)


func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}
	
    // 1. Rodar Migrations primeiro (Usa database/sql por baixo)
	if err := database.RunMigrations(dbURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
    }
	
	ctx := context.Background()

	// 2. Iniciar Pool de Conexão para a aplicação (Usa pgxpool)
    pool, err := database.NewPool(ctx, dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer pool.Close()

    // 3. Iniciar Servidor gRPC passando o pool como dependência
    // server := grpc.NewServer(...)
    // accounts.RegisterAccountServiceServer(server, &handler.AccountHandler{DB: pool})
}