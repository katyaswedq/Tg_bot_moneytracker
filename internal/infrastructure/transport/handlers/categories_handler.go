package handlers

import (
	"context"
	"strings"

	"tgfin/internal/infrastructure/transport/transportdto"
	"tgfin/internal/service/contract"
	"tgfin/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CategoriesHandler struct {
	categoryService contract.CategoryService
}

func NewCategoriesHandler(categoryService contract.CategoryService) *CategoriesHandler {
	return &CategoriesHandler{
		categoryService: categoryService,
	}
}

func (h *CategoriesHandler) Handle(ctx context.Context, update tgbotapi.Update) (*transportdto.Result, error) {
	if update.Message == nil || update.Message.From == nil {
		return nil, nil
	}

	out, err := h.categoryService.List(ctx, dto.CategoryListInput{
		TelegramID: update.Message.From.ID,
	})
	if err != nil {
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Произошла ошибка. Попробуйте позже.",
		}, err
	}

	if len(out.Items) == 0 {
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "📂 У вас пока нет категорий.\nИспользуйте /category add <название>",
		}, nil
	}

	var b strings.Builder
	b.WriteString("📂 Ваши категории:\n\n")

	for _, c := range out.Items {
		b.WriteString("• ")
		b.WriteString(c.Name)
		b.WriteString("\n   ID: ")
		b.WriteString(c.ID.String())
		b.WriteString("\n\n")
	}

	b.WriteString("💡 Используйте ID для удаления категории\n")
	b.WriteString("Например: /category delete <id>")

	return &transportdto.Result{
		ChatID: update.Message.Chat.ID,
		Text:   b.String(),
	}, nil
}
