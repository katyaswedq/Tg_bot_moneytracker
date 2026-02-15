package dto

import "github.com/google/uuid"

type CategoryAddInput struct {
	TelegramID int64
	Name       string
}

type CategoryAddOutput struct {
	ID   uuid.UUID
	Name string
}

type CategoryListInput struct {
	TelegramID int64
}

type CategoryListItem struct {
	ID   uuid.UUID
	Name string
}

type CategoryListOutput struct {
	Items []CategoryListItem
}

type CategoryDeleteInput struct {
	TelegramID int64
	CategoryID uuid.UUID
}
