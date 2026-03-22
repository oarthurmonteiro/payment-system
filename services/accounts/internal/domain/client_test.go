package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("should create a valid client with success", func(t *testing.T) {
		// Dados de entrada válidos
		name := "Arthur Monteiro"
		docRaw := "12345678909" // CPF que passe no NewDocument

		client, err := domain.NewClient(name, docRaw)

		// Asserções
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, name, client.FullName)
		assert.Equal(t, docRaw, client.Document.Value())
		
		// Valida se o ID foi gerado (V7 é superior a zero)
		assert.NotEqual(t, uuid.Nil, client.ID)
		
		// Valida se o timestamp foi criado recentemente (tolerância de 1s)
		assert.WithinDuration(t, time.Now(), client.CreatedAt, time.Second)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		client, err := domain.NewClient("", "12345678909")

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "Full name is required")
	})

	t.Run("should return error when document is invalid", func(t *testing.T) {
		// NewClient falha porque o NewDocument interno falha
		client, err := domain.NewClient("Arthur", "123") 

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestRestoreClient(t *testing.T) {
	t.Run("should restore an existing client from database data", func(t *testing.T) {
		// Dados simulando o que veio do Postgres
		id, _ := uuid.NewV7()
		name := "Arthur Monteiro"
		doc := "12345678909"
		createdAt := time.Now().Add(-24 * time.Hour) // Ontem

		client := domain.RestoreClient(id, name, doc, createdAt)

		assert.Equal(t, id, client.ID)
		assert.Equal(t, name, client.FullName)
		assert.Equal(t, doc, client.Document.Value())
		assert.Equal(t, createdAt, client.CreatedAt)
	})
}