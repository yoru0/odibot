package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yoru0/goodi/internal/core/domain"
	"github.com/yoru0/goodi/internal/ports"
)

type calculatorService struct {
	calcRepo  ports.CalculationRepository
	cacheRepo ports.CacheRepository
}

func NewCalculatorService(calcRepo ports.CalculationRepository, cacheRepo ports.CacheRepository) ports.CalculatorService {
	return &calculatorService{
		calcRepo:  calcRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *calculatorService) Add(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error) {
	return s.performCalculation(ctx, userID, "add", a, b, a+b)
}

func (s *calculatorService) Subtract(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error) {
	return s.performCalculation(ctx, userID, "subtract", a, b, a-b)
}

func (s *calculatorService) Multiply(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error) {
	return s.performCalculation(ctx, userID, "multiply", a, b, a*b)
}

func (s *calculatorService) Divide(ctx context.Context, userID string, a, b float64) (*domain.CalculationResult, error) {
	if b == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	return s.performCalculation(ctx, userID, "divide", a, b, a/b)
}

func (s *calculatorService) performCalculation(
	ctx context.Context,
	userID, operation string,
	a, b, result float64,
) (*domain.CalculationResult, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("calc:%s:%.2f:%.2f", operation, a, b)
	if cachedResult, err := s.cacheRepo.Get(ctx, cacheKey); err == nil && cachedResult != "" {
		// Return cached result without saving to DB again
		return &domain.CalculationResult{
			Result: result,
			Calculation: &domain.Calculation{
				UserID:    userID,
				Operation: operation,
				Operand1:  a,
				Operand2:  b,
				Result:    result,
			},
		}, nil
	}

	// Create calculation record
	calc := &domain.Calculation{
		ID:        uuid.New().String(),
		UserID:    userID,
		Operation: operation,
		Operand1:  a,
		Operand2:  b,
		Result:    result,
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := s.calcRepo.SaveCalculation(ctx, calc); err != nil {
		return nil, fmt.Errorf("failed to save calculation: %w", err)
	}

	// Cache the result for 1 hour
	s.cacheRepo.Set(ctx, cacheKey, fmt.Sprintf("%.2f", result), 3600)

	return &domain.CalculationResult{
		Result:      result,
		Calculation: calc,
	}, nil
}