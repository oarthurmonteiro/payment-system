package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Erros específicos de domínio
var (
	ErrRequiredDocument = errors.New("Document is required")
	ErrInvalidDocument  = errors.New("Document invalid")
)

type Document struct {
    value string
}

func NewDocument(v string) (Document, error) {
	cleaned := clean(v)

	if strings.TrimSpace(cleaned) == "" {
		return Document{}, ErrRequiredDocument
	}

	if err := validateCPF(cleaned); err != nil {
		return Document{}, err
	}

	return Document{value: cleaned}, nil
}

func RestoreDocument(v string) Document {
	return Document{value: v}
}

func (d Document) Value() string {
	return d.value
}

func validateCPF(v string) error {
	
	if len(v) != 11 {
		return fmt.Errorf("%w: must have 11 digits", ErrInvalidDocument)
	}

	if hasAllSameDigits(v) {
		return fmt.Errorf("%w: same digits", ErrInvalidDocument)
	}

	if !remaindersAreValid(v) {
		return ErrInvalidDocument
	}

	return nil
}

func clean(input string) string {
	var sb strings.Builder
	for _, r := range input {
		if unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func hasAllSameDigits(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func remaindersAreValid(cpf string) bool {
	var sum int
	multipliers1 := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	
	for i := range 9 {
		sum += int(cpf[i]-'0') * multipliers1[i]
	}
	
	remainder := sum % 11
	if remainder < 2 {
		remainder = 0
	} else {
		remainder = 11 - remainder
	}

	if int(cpf[9]-'0') != remainder {
		return false
	}

	sum = 0
	multipliers2 := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
	for i := range 10 {
		sum += int(cpf[i]-'0') * multipliers2[i]
	}

	remainder = sum % 11
	if remainder < 2 {
		remainder = 0
	} else {
		remainder = 11 - remainder
	}

	return int(cpf[10]-'0') == remainder
}