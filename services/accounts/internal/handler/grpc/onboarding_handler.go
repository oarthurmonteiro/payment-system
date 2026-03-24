package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/oarthurmonteiro/payment-system/services/accounts/gen/accounts/v1"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/domain"
	"github.com/oarthurmonteiro/payment-system/services/accounts/internal/usecase"
)

type OnboardingUseCase interface {
	Execute(ctx context.Context, input usecase.OnboardingInput) (usecase.OnboardingOutput, error)
}

type OnboardingHandler struct {
	// O gRPC exige que você embarque essa struct para compatibilidade futura
	pb.UnimplementedOnboardingServiceServer
	useCase OnboardingUseCase
}

func NewOnboardingHandler(uc OnboardingUseCase) *OnboardingHandler {
	return &OnboardingHandler{useCase: uc}
}

func (h *OnboardingHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 1. Chamar o UseCase (Lógica de Negócio)
	// Criamos um DTO de input para não passar o "req" do gRPC direto para o core
	input := usecase.OnboardingInput{
		Name: req.GetFullName(),
		Document:  req.GetDocument(),
	}

	output, err := h.useCase.Execute(ctx, input)
	if err != nil {
		return nil, h.mapError(err)
	}

	// 2. Mapear o Output do UseCase para a Resposta do gRPC (Protobuf)
	return &pb.RegisterResponse{
		ClientId:  output.ClientID.String(),
		AccountId: output.AccountID.String(),
		Status:    pb.AccountStatus_ACCOUNT_STATUS_PENDING,
		CreatedAt: timestamppb.New(output.CreatedAt),
	}, nil
}


// mapError traduz erros de negócio para códigos de resposta gRPC
func (h *OnboardingHandler) mapError(err error) error {
	// 1. Erro de validação de CPF (400 Bad Request no REST -> InvalidArgument no gRPC)
    if errors.Is(err, domain.ErrRequiredDocument) || errors.Is(err, domain.ErrInvalidDocument) {
        return status.Error(codes.InvalidArgument, err.Error())
    }

    // 2. Erro de cliente já existente (409 Conflict -> AlreadyExists no gRPC)
    // Se você tiver essa verificação no Use Case
    if errors.Is(err, domain.ErrClientAlreadyExists) {
        return status.Error(codes.AlreadyExists, "this document is already registered")
    }

	if errors.Is(err, domain.ErrEventSerialization) {
        // Logamos o erro aqui para o desenvolvedor ver no CloudWatch/Grafana
        log.Printf("Critico: erro de outbox: %v", err)
        return status.Error(codes.Internal, "erro interno ao processar evento")
    }

    // 3. Fallback para erros inesperados (500 Internal Server Error)
    // Logamos o erro real aqui antes de esconder do cliente por segurança
	log.Printf("Erro não mapeado no Onboarding: %v", err)
    return status.Error(codes.Internal, "internal server error")
}