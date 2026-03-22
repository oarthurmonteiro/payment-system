package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/usecase"
)

type mockRepo struct {
	err error
}

func (m *mockRepo) CreateOnboarding(ctx context.Context, c *domain.Client, a *domain.Account, e *domain.OutboxEvent) error {
	return m.err
}

func TestOnboardingUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		input         usecase.OnboardingInput
		repoError     error
		expectedError bool
	}{
		{
			name: "Should onboarding successfully",
			input: usecase.OnboardingInput{
				Name:     "Arthur Monteiro",
				Document: "30011013770", // Aquele CPF que descobrimos que o resto é < 2
			},
			repoError:     nil,
			expectedError: false,
		},
		{
			name: "Should return error if CPF is invalid",
			input: usecase.OnboardingInput{
				Name:     "Arthur",
				Document: "123", // Inválido, vai barrar no domain.NewClient
			},
			repoError:     nil,
			expectedError: true,
		},
		{
			name: "Should return error if repository fails",
			input: usecase.OnboardingInput{
				Name:     "Arthur",
				Document: "30011013770",
			},
			repoError:     errors.New("db connection failed"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := &mockRepo{err: tt.repoError}
			uc := usecase.NewOnboardingUseCase(repo)

			// Act
			output, err := uc.Execute(context.Background(), tt.input)

			// Assert
			if (err != nil) != tt.expectedError {
				t.Errorf("Execute() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if !tt.expectedError {
				if output.ClientID == uuid.Nil || output.AccountID == uuid.Nil {
					t.Error("Execute() returned empty UUIDs on success")
				}
				if output.Status == "" {
					t.Error("Execute() returned empty status")
				}
			}
		})
	}
}