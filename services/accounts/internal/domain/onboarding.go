package domain

import "context"

type OnboardingRepository interface {
    // CreateOnboarding executa a transação atômica: Client + Account + Outbox
    CreateOnboarding(ctx context.Context, c *Client, a *Account, e *OutboxEvent) error
}