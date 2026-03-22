package database_test

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/database"
)

type PostgresSuiteBase struct {
    suite.Suite
    container *postgres.PostgresContainer
    pool      *pgxpool.Pool
    repo      *database.PostgresRepository
    ctx       context.Context
}

// SetupSuite roda UMA VEZ antes de todos os testes da struct (Sobe o Docker)
func (s *PostgresSuiteBase) SetupSuite() {
    s.ctx = context.Background()
    
    container, err := postgres.RunContainer(s.ctx,
        testcontainers.WithImage("postgres:18-alpine"),
		postgres.WithDatabase("accounts_test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
        testcontainers.WithWaitStrategy(
            wait.
                ForLog("database system is ready").
                WithOccurrence(2)),
    )
    s.NoError(err)
    s.container = container

    connStr, _ := container.ConnectionString(s.ctx, "sslmode=disable")
    
    // Roda as migrations (seu iofs)
    err = database.RunMigrations(connStr)
    s.NoError(err)

    pool, _ := database.NewPostgresPool(s.ctx, connStr)
    s.pool = pool
    s.repo = database.NewPostgresRepository(pool)
}

// TearDownSuite roda UMA VEZ após todos os testes (Fecha o Docker)
func (s *PostgresSuiteBase) TearDownSuite() {
    s.pool.Close()
    s.container.Terminate(s.ctx)
}

// SetupTest RODA ANTES DE CADA TESTE (A mágica da limpeza acontece aqui!)
func (s *PostgresSuiteBase) SetupTest() {
    err := cleanDatabase(s.ctx, s.pool)
    s.NoError(err)
}

func cleanDatabase(ctx context.Context, pool *pgxpool.Pool) error {
    tables := []string{
        "outbox",
        "accounts",
        "clients",
    }

    for _, table := range tables {
        // TRUNCATE com CASCADE limpa as FKs relacionadas e RESTART IDENTITY reseta IDs seriais
        _, err := pool.Exec(ctx, "TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE")
        if err != nil {
            return err
        }
    }
    return nil
}