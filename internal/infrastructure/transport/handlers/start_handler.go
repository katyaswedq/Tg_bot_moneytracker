package handlers

import (
	"context"

	"tgfin/internal/infrastructure/transport/transportdto"
	"tgfin/internal/service/contract"
	"tgfin/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	startService contract.StartService
}

func NewStartHandler(startService contract.StartService) *StartHandler {
	return &StartHandler{
		startService: startService,
	}
}

func (h *StartHandler) Handle(ctx context.Context, update tgbotapi.Update) (*transportdto.Result, error) {
	if update.Message == nil || update.Message.From == nil {
		return nil, nil
	}

	from := update.Message.From

	var username *string
	if from.UserName != "" {
		u := from.UserName
		username = &u
	}

	in := dto.StartInput{
		TelegramID: from.ID,
		UserName:   username,
		FirstName:  from.FirstName,
	}

	if err := h.startService.Start(ctx, in); err != nil {
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Произошла ошибка. Попробуйте позже.",
		}, err
	}

	return &transportdto.Result{
		ChatID: update.Message.Chat.ID,
		Text: "👋 Добро пожаловать!\n\n" +
			"Я помогу вам отслеживать расходы и управлять бюджетами.\n\n" +
			"✅ Вы зарегистрированы!\n" +
			"📂 Созданы базовые категории:\n" +
			"   • Еда\n" +
			"   • Транспорт\n" +
			"   • Развлечения\n" +
			"   • Прочее\n\n" +
			"Используйте /help для списка команд",
	}, nil
}
