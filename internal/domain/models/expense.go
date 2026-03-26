package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	CategoryID  uuid.UUID
	Amount      int64
	Description *string
	CreatedAt   time.Time
}

func NewExpense(userID, categoryID uuid.UUID, amount int64, description *string) *Expense {
	return &Expense{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      amount,
		Description: description,
	}
}
