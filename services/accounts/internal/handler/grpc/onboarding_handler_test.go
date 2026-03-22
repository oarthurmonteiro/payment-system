package grpc_test

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/oarthurmonteiro/payment-system/services/accounts/gen/accounts/v1"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/handler/grpc"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/usecase"
)

// mockRepo implementa a interface domain.OnboardingRepository
type mockRepo struct {
    err error
}

func (m *mockRepo) CreateOnboarding(ctx context.Context, c *domain.Client, a *domain.Account, e *domain.OutboxEvent) error {
    return m.err
}

func TestOnboardingHandler_Register(t *testing.T) {
    t.Run("Should return InvalidArgument when CPF is invalid", func(t *testing.T) {
        // Para erro de CPF, o repositório nem chega a ser chamado, então o mock pode ser vazio
        repo := &mockRepo{}
        uc := usecase.NewOnboardingUseCase(repo)
        handler := grpc.NewOnboardingHandler(uc)

        req := &pb.RegisterRequest{
            FullName: "Arthur Monteiro",
            Document: "123", // Inválido
        }

        _, err := handler.Register(context.Background(), req)

        st, ok := status.FromError(err)
        if !ok || st.Code() != codes.InvalidArgument {
            t.Errorf("expected InvalidArgument, got %v", st.Code())
        }
    })

    t.Run("Should return Success when data is valid", func(t *testing.T) {
        repo := &mockRepo{err: nil} // Sucesso no banco
        uc := usecase.NewOnboardingUseCase(repo)
        handler := grpc.NewOnboardingHandler(uc)

        req := &pb.RegisterRequest{
            FullName: "Arthur Monteiro",
            Document: "30011013770", 
        }

        resp, err := handler.Register(context.Background(), req)

        if err != nil {
            t.Fatalf("expected no error, got %v", err)
        }

        if resp.GetClientId() == "" || resp.GetAccountId() == "" {
            t.Error("expected IDs in response")
        }
    })

    t.Run("Should return AlreadyExists when Client already exists", func(t *testing.T) {
        // CONFIGURAÇÃO: O mock vai forçar o erro de duplicidade
        repo := &mockRepo{err: domain.ErrClientAlreadyExists}
        uc := usecase.NewOnboardingUseCase(repo)
        handler := grpc.NewOnboardingHandler(uc)

        req := &pb.RegisterRequest{
            FullName: "Arthur Monteiro",
            Document: "30011013770",
        }

        // AÇÃO: Chamamos o Handler (como um cliente gRPC faria)
        _, err := handler.Register(context.Background(), req)

        // VERIFICAÇÃO: O seu mapError deve ter convertido para AlreadyExists
        st, ok := status.FromError(err)
        if !ok {
            t.Fatal("expected gRPC status error")
        }

        if st.Code() != codes.AlreadyExists {
            t.Errorf("expected AlreadyExists, got %v", st.Code())
        }
        
        expectedMsg := "this document is already registered"
        if st.Message() != expectedMsg {
            t.Errorf("expected message %s, got %s", expectedMsg, st.Message())
        }
    })

	t.Run("Should return Internal and log critical when event serialization fails", func(t *testing.T) {
		// Simulamos que o erro aconteceu "dentro" do processo
		repo := &mockRepo{err: errors.New("Some other error")}
		uc := usecase.NewOnboardingUseCase(repo)
		handler := grpc.NewOnboardingHandler(uc)

		req := &pb.RegisterRequest{
			FullName: "Arthur Monteiro",
			Document: "30011013770",
		}

		_, err := handler.Register(context.Background(), req)

		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Internal {
			t.Fatalf("expected Internal error, got %v", st.Code())
		}
	})
}