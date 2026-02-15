package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	Name      string
	Color     *string
	CreatedAt time.Time
	IsDefault bool
}

func NewCategory(userId uuid.UUID, name string, isDefault bool) *Category {
	return &Category{
		UserId:    userId,
		Name:      name,
		CreatedAt: time.Now(),
		IsDefault: isDefault,
	}
}
