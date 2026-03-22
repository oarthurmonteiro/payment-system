package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	t.Run("should create a valid outbox event with serialized payload", func(t *testing.T) {
		// Arrange: Criamos um dado qualquer para simular o evento
		type SampleData struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		data := SampleData{ID: uuid.NewString(), Name: "Arthur"}
		eventType := "client.created"

		// Act
		event, err := domain.NewEvent(eventType, data)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, event)
		assert.NotEqual(t, uuid.Nil, event.ID)
		assert.Equal(t, eventType, event.EventType)
		assert.Equal(t, domain.OutboxStatusPending, event.Status)
		
		// 1. Validar se o payload é um JSON válido e contém os dados originais
		var decoded SampleData
		err = json.Unmarshal(event.Payload, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, data.Name, decoded.Name)
	})

	t.Run("should return error when data cannot be serialized", func(t *testing.T) {
		// No Go, canais (chan) ou funções não podem ser serializados para JSON
		invalidData := make(chan int)

		event, err := domain.NewEvent("test.event", invalidData)

		assert.ErrorIs(t, err, domain.ErrEventSerialization)
		assert.Nil(t, event)
	})
}