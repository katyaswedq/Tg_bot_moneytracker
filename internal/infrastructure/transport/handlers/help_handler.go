package handlers

import (
	"context"

	"tgfin/internal/infrastructure/transport/transportdto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct{}

func NewHelpHandler() *HelpHandler {
	return &HelpHandler{}
}

func (h *HelpHandler) Handle(ctx context.Context, update tgbotapi.Update) (*transportdto.Result, error) {
	if update.Message == nil {
		return nil, nil
	}

	text :=
		"📖 Доступные команды:\n\n" +
			"💰 Расходы:\n" +
			"/add <сумма> <категория> <описание> — добавить расход\n\n" +
			"📂 Категории:\n" +
			"/category add <название> — создать категорию\n" +
			"/categories — список категорий\n" +
			"/category delete <id> — удалить категорию\n\n" +
			"ℹ️ Прочее:\n" +
			"/start — регистрация\n" +
			"/help — справка"

	return &transportdto.Result{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	}, nil
}
