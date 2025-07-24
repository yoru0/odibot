package ports

import (
	"context"

	"github.com/yoru0/goodi/internal/core/domain"
)

type CalculationRepository interface {
	SaveCalculation(ctx context.Context, calc *domain.Calculation) error
	GetCalculationsByUser(ctx context.Context, userID string, limit int) ([]*domain.Calculation, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}