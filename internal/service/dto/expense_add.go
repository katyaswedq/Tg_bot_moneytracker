package dto

import "github.com/google/uuid"

type ExpenseAddInput struct {
	TelegramID   int64
	Amount       int64
	CategoryName string
	Description  *string
}

type ExpenseAddOutput struct {
	ID          uuid.UUID
	Amount      int64
	Category    string
	Description *string
}
