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
		err		  error
	}{
		{
			name:      	"Should accept a valid CPF",
			input:     	"12345678909",
			shouldErr: 	false,
			err:	   	nil,
		},
		{
			name:      	"Should accept CPF where 2nd digit remainder < 2 (Result 0)",
			input:     	"30011013770",
			shouldErr: 	false,
			err:	    nil,
		},
		{
			name:     	"Should reject CPF with invalid first digit",
			input:    	"45443353801", // Should be 8
			shouldErr: 	true,
			err:		domain.ErrInvalidDocument,
		},
		{
			name:     	"Should reject CPF with invalid second digit",
			input:    	"30011013771", // Should be 0
			shouldErr: 	true,
			err:		domain.ErrInvalidDocument,
		},
		{
			name:      	"Should reject incomplete Document",
			input:     	"123",
			shouldErr: 	true,
			err:		domain.ErrInvalidDocument,
		},
		{
			name:	   	"Should reject Document with all digits the same",
			input: 	   	"111.111.111-11",
			shouldErr: 	true,
			err:		domain.ErrInvalidDocument,
		},
		{
			name:      	"Should reject Document with letters",
			input:     	"1234567890A",
			shouldErr: 	true,
			err:		domain.ErrInvalidDocument,
		},
		{
			name:      	"Should reject empty Document",
			input:     	"",
			shouldErr: 	true,
			err:		domain.ErrRequiredDocument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := domain.NewDocument(tt.input)

			if tt.shouldErr {
				assert.ErrorIs(t, err, tt.err)
				assert.Equal(t, domain.Document{}, doc, "Para erro, a struct deve vir zerada")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, doc)
				assert.Equal(t, tt.input, doc.Value())
			}
		})
	}
}