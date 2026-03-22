package domain_test

import (
	"testing"

	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewDocument(t *testing.T) {
	// Definimos os cenários de teste em uma struct anônima
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "Deve aceitar um CPF válido",
			input:     "12345678909",
			shouldErr: false,
		},
		{
			name:      "Deve rejeitar um CPF muito curto",
			input:     "123",
			shouldErr: true,
		},
		{
			name:      "Deve rejeitar um CPF com letras",
			input:     "1234567890A",
			shouldErr: true,
		},
		{
			name:      "Deve rejeitar um campo vazio",
			input:     "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Chamamos a função construtora do seu domínio
			doc, err := domain.NewDocument(tt.input)

			if tt.shouldErr {
				assert.Error(t, err, "Esperava erro para o input: %s", tt.input)
				assert.Equal(t, domain.Document{}, doc, "Para erro, a struct deve vir zerada")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, doc)
				assert.Equal(t, tt.input, doc.Value())
			}
		})
	}
}