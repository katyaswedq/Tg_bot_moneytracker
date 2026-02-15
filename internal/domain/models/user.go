package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID         uuid.UUID
    TelegramID int64
    UserName   *string
    FirstName  string
    CreatedAt  time.Time
}

func NewUser(telegramID int64, userName *string, firstName string) *User {
    return &User{
        TelegramID: telegramID,
        UserName:   userName,
        FirstName:  firstName,
        CreatedAt:  time.Now(),
    }
}