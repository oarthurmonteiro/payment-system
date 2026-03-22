package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrEventSerialization = errors.New("falha ao serializar dados do evento")

type OutboxEvent struct {
	ID        uuid.UUID
	EventType string
	Payload   []byte // JSON serializado
	Status    OutboxStatus
	CreatedAt time.Time
}

type OutboxStatus string
const (
	OutboxStatusPending OutboxStatus = "PENDING"
	OutboxStatusProcessed OutboxStatus = "PROCESSED"
	OutboxStatusFailed OutboxStatus = "FAILED"
)

func (s OutboxStatus) String() string {
	return string(s)
}

func NewEvent(eventType string, data any) (*OutboxEvent, error) {
	// 1. Serializa o objeto de dados para JSON (armazenado como []byte)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, ErrEventSerialization
	}

	// 2. Retorna o ponteiro para a Entidade OutboxEvent
	return &OutboxEvent{
		ID:        uuid.New(),           // ID único da mensagem
		EventType: eventType,            // Ex: "account.created"
		Payload:   payload,
		Status:    OutboxStatusPending,  // Sempre começa como PENDING
		CreatedAt: time.Now(),
	}, nil
}

