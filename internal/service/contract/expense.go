package contract

import (
	"context"

	"tgfin/internal/service/dto"
)

type ExpenseService interface {
	Add(ctx context.Context, in dto.ExpenseAddInput) (*dto.ExpenseAddOutput, error)
}
