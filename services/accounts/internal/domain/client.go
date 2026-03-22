package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrClientAlreadyExists = errors.New("Client already in the system")

type Client struct {
	ID        uuid.UUID
	FullName  string
	Document  Document
	CreatedAt time.Time
}

func NewClient(name string, docRaw string) (*Client, error) {
	// 1. Valida e cria o Value Object primeiro
	document, err := NewDocument(docRaw)
	if err != nil {
		return nil, err // Retorna o erro envelopado (%w) que definimos
	}

	// 2. Valida regras de nome (ex: não pode ser vazio)
	if name == "" {
		return nil, errors.New("Full name is required.")
	}

	// 3. Gera os dados de sistema (ID e Timestamp)
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Client{
		ID:        id,
		FullName:  name,
		Document:  document,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// RestoreClient é usado pelo Repository para transformar linhas do SQL em objetos de Domínio
func RestoreClient(id uuid.UUID, name string, document string, createdAt time.Time) *Client {
	return &Client{
		ID:        id,
		FullName:  name,
		Document:  RestoreDocument(document),
		CreatedAt: createdAt,
	}
}