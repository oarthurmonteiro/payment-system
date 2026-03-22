package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
)

type OnboardingInput struct {
    Name string
    Document  string
}

// OnboardingOutput também deve ser exportada para o Handler ler os resultados
type OnboardingOutput struct {
    ClientID  uuid.UUID
    AccountID uuid.UUID
	Status    string
    CreatedAt time.Time
}

type OnboardingUseCase struct {
    repo domain.OnboardingRepository
}

func NewOnboardingUseCase(r domain.OnboardingRepository) *OnboardingUseCase {
    return &OnboardingUseCase{repo: r}
}

func (uc *OnboardingUseCase) Execute(ctx context.Context, input OnboardingInput) (OnboardingOutput, error) {
    // 1. Cria o Cliente (Domínio)
    client, err := domain.NewClient(input.Name, input.Document)
    if err != nil {
        return OnboardingOutput{}, err
    }

    // 2. Cria a Conta (Domínio)
    account, err := domain.NewAccount(client.ID)
	if err != nil {
		return OnboardingOutput{}, err
	}

	eventData := map[string]interface{}{
        "account_id": account.ID,
        "client_id":  client.ID,
        "document":   client.Document.Value(),
    }
	event, err := domain.NewEvent("account.created", eventData)
    if err != nil {
        return OnboardingOutput{}, err
    }

	if err := uc.repo.CreateOnboarding(ctx, client, account, event); err != nil {
        return OnboardingOutput{}, err
    }

    // 3. Persiste (Repository)
    return OnboardingOutput{
        ClientID:  client.ID,
        AccountID: account.ID,
        Status:    account.Status.String(),
        CreatedAt: account.CreatedAt,
    }, nil
}