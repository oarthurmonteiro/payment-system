package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrRequiredClientID = errors.New("Client ID is required")

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

func NewAccount(clientID uuid.UUID) (*Account, error) {
	if clientID == uuid.Nil {
        return nil, ErrRequiredClientID
    }
	
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Account{
		ID: 	  id,
		ClientID: clientID,
		Status:  AccountStatusPending,
		CreatedAt: time.Now().UTC(),
	}, nil
}