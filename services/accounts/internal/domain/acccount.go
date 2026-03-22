package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID
	ClientID  uuid.UUID
	Status    AccountStatus
	CreatedAt time.Time
}

type AccountStatus string

const (
    AccountStatusPending  AccountStatus = "PENDING"
    AccountStatusActive   AccountStatus = "ACTIVE"
    AccountStatusBlocked  AccountStatus = "BLOCKED"
    AccountStatusCanceled AccountStatus = "CANCELED"
)

// 3. (Opcional) Método para facilitar o uso ou validação
func (s AccountStatus) String() string {
    return string(s)
}

func NewAccount(clientId uuid.UUID) (*Account, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Account{
		ID: 	  id,
		ClientID: clientId,
		Status:  AccountStatusPending,
		CreatedAt: time.Now().UTC(),
	}, nil
}