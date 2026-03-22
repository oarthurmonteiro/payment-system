package database_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/stretchr/testify/suite"
)

func TestOnboardingRepoSuite(t *testing.T) {
    suite.Run(t, new(PostgresSuiteBase))
}

func (s *PostgresSuiteBase) TestCreateOnboarding_Success() {
    client, _ := domain.NewClient("Arthur Monteiro", "30011013770")
    account, _ := domain.NewAccount(client.ID)
    event := &domain.OutboxEvent{
            ID: uuid.New(),
            EventType: "account.created",
            Payload: []byte(`{"test":true}`),
            Status: domain.OutboxStatusPending,
            CreatedAt: time.Now(),
        }

    err := s.repo.CreateOnboarding(s.ctx, client, account, event)

    s.NoError(err)
    
    var exists bool
	err = s.pool.QueryRow(s.ctx, "SELECT EXISTS(SELECT 1 FROM clients WHERE id = $1)", client.ID).Scan(&exists)
	s.NoError(err)
    s.True(exists)
}

func (s *PostgresSuiteBase) TestCreateOnboarding_Duplicate() {
    client, _ := domain.NewClient("Arthur", "30011013770")
	account, _ := domain.NewAccount(client.ID)
	event := &domain.OutboxEvent{
            ID: uuid.New(),
            EventType: "account.created",
            Payload: []byte(`{"test":true}`),
            Status: domain.OutboxStatusPending,
            CreatedAt: time.Now(),
        }
    err := s.repo.CreateOnboarding(s.ctx, client, account, event)
	s.NoError(err)

    // Segunda tentativa com mesmo CPF
	client_other, _ := domain.NewClient("Arthur", "30011013770")
	account_other, _ := domain.NewAccount(client_other.ID)
	event_other := &domain.OutboxEvent{
            ID: uuid.New(),
            EventType: "account.created",
            Payload: []byte(`{"test":true}`),
            Status: domain.OutboxStatusPending,
            CreatedAt: time.Now(),
        }
    err = s.repo.CreateOnboarding(s.ctx, client_other, account_other, event_other)
    s.ErrorIs(err, domain.ErrClientAlreadyExists)
}