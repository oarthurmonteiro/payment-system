package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
)

type PostgresRepository struct {
    pool *pgxpool.Pool
}

// NewPool cria um pool de conexões com configurações de performance
func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, err
    }

    // Configurações de "Stress Test"
    config.MaxConns = 100
    config.MinConns = 10
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = time.Minute * 30

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, err
    }

    // Ping para garantir que o banco está vivo
    if err := pool.Ping(ctx); err != nil {
        return nil, err
    }

    return pool, nil
}

func (r *PostgresRepository) handlePgError(err error) error {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        if pgErr.Code == "23505" { // Unique Violation
            // O segredo está aqui: olhar QUAL constraint falhou
            switch pgErr.ConstraintName {
            case "clients_document_key":
                return domain.ErrClientAlreadyExists
            // case "accounts_number_key":
            //     return domain.ErrAccountNumberAlreadyExists // Novo erro de domínio
            // case "idx_unique_active_account":
            //     return domain.ErrOnlyOneActiveAccountAllowed
            }
        }
    }
    return err
}