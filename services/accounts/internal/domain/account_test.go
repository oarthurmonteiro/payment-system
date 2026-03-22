package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	t.Run("should create a valid account with success", func (t *testing.T) {
		// Dados de entrada válidos
		clientId, _ := uuid.NewV7()

		account, err := domain.NewAccount(clientId)

		// Asserções
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, clientId, account.ClientID)

		assert.NotEqual(t, uuid.Nil, account.ID)

		assert.WithinDuration(t, time.Now(), account.CreatedAt, time.Second)

	})

	t.Run("should return error when clientId is empty", func (t *testing.T) {
		account, err := domain.NewAccount(uuid.Nil)

		assert.Nil(t, account)
		assert.ErrorIs(t, err, domain.ErrRequiredClientID)
	})
}