package ports

import (
	"context"

	"github.com/yoru0/goodi/internal/core/domain"
)

type CalculatorService interface {
	Add(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error)
	Subtract(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error)
	Multiply(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error)
	Divide(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error)
}

type CommandHandler interface {
	Handle(ctx context.Context, userID, channelID string, args []string) (string, error)
	GetUsage() string
}