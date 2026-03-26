package user

import (
	"context"

	"github.com/google/uuid"
	"tgfin/internal/domain/models"
)

type Tx interface{}

type CommandTag interface {
	RowsAffected() int64
}

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

type Exec interface {
	Exec(ctx context.Context, sql string, args ...any) error
	ExecWithResult(ctx context.Context, sql string, args ...any) (CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

type UnitOfWork interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, exec Exec) error) error
}

type UserRepository interface {
	Create(ctx context.Context, exec Exec, user *models.User) error
	GetByTelegramID(ctx context.Context, exec Exec, telegramID int64) (*models.User, error)
}

type CategoryRepository interface {
	CreateDefault(ctx context.Context, exec Exec, userID uuid.UUID) error
	Create(ctx context.Context, exec Exec, userID uuid.UUID, name string) (*models.Category, error)
	ListByUser(ctx context.Context, exec Exec, userID uuid.UUID) ([]*models.Category, error)
	Delete(ctx context.Context, exec Exec, userID uuid.UUID, categoryID uuid.UUID) error
	GetByName(ctx context.Context, exec Exec, userID uuid.UUID, name string) (*models.Category, error)
}

type ExpenseRepository interface {
	Create(ctx context.Context, exec Exec, e *models.Expense) error
}