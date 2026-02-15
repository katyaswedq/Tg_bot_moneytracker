package app

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Task struct {
	Update tgbotapi.Update
}

type Delivery struct {
	ChatID int64
	Text   string
}
