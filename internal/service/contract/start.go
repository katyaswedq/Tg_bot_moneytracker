package contract

import (
	"context"

	"tgfin/internal/service/dto"
)

type StartService interface {
	Start(ctx context.Context, in dto.StartInput) error
}
