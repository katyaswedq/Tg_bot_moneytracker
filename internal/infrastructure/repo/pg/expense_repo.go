package pg

import (
	"context"

	"tgfin/internal/domain/models"
	"tgfin/internal/domain/usecase"
)

type ExpenseRepo struct{}

func NewExpenseRepo() *ExpenseRepo {
	return &ExpenseRepo{}
}

func (r *ExpenseRepo) Create(ctx context.Context, exec user.Exec, e *models.Expense) error {
	query := `
		INSERT INTO expenses (user_id, category_id, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return exec.QueryRow(ctx, query, e.UserID, e.CategoryID, e.Amount,e.Description).Scan(&e.ID, &e.CreatedAt)
}
