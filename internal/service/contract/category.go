package contract

import (
	"context"

	"tgfin/internal/service/dto"
)

type CategoryService interface {
	Add(ctx context.Context, in dto.CategoryAddInput) (*dto.CategoryAddOutput, error)
	List(ctx context.Context, in dto.CategoryListInput) (*dto.CategoryListOutput, error)
	Delete(ctx context.Context, in dto.CategoryDeleteInput) error
}
