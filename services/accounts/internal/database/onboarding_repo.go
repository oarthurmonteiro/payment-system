package database

import (
	"context"

	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
)

func (r *PostgresRepository) CreateOnboarding(
    ctx context.Context, 
    c *domain.Client, 
    a *domain.Account,
    e *domain.OutboxEvent,
) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return err
    }
	
    // O defer garante que nada seja salvo se houver erro no meio
    defer tx.Rollback(ctx)

    if _, err := tx.Exec(ctx, `
        INSERT INTO clients (id, full_name, document, created_at)
        VALUES ($1, $2, $3, $4)
    `, c.ID, c.FullName, c.Document.Value(), c.CreatedAt); err != nil {
        return r.handlePgError(err)
    }
    
    if _, err := tx.Exec(ctx, `
        INSERT INTO accounts (id, client_id, status, created_at)
        VALUES ($1, $2, $3, $4)
    `, a.ID, a.ClientID, a.Status, a.CreatedAt); err != nil {
        return r.handlePgError(err)
    }

    if _, err := tx.Exec(ctx, `
        INSERT INTO outbox (id, topic, payload, status, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `, e.ID, e.EventType, e.Payload, e.Status.Value(), e.CreatedAt); err != nil { 
        return r.handlePgError(err) 
    }

    return tx.Commit(ctx)
}