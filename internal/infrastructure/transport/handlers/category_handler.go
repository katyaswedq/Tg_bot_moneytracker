package handlers

import (
	"context"
	"strings"

	"github.com/google/uuid"

	domainerr "tgfin/internal/domain/error"
	"tgfin/internal/infrastructure/transport/transportdto"
	"tgfin/internal/service/contract"
	"tgfin/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CategoryHandler struct {
	categoryService contract.CategoryService
}

func NewCategoryHandler(categoryService contract.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Handle(ctx context.Context, update tgbotapi.Update) (*transportdto.Result, error) {
	if update.Message == nil || update.Message.From == nil {
		return nil, nil
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "📂 Использование:\n/category add <название>\n/category delete <id>",
		}, nil
	}

	parts := strings.Fields(args)
	sub := strings.ToLower(parts[0])

	switch sub {
	case "add":
		name := strings.TrimSpace(strings.TrimPrefix(args, parts[0]))
		if name == "" {
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Укажите название.\nПример: /category add Спорт",
			}, nil
		}

		out, err := h.categoryService.Add(ctx, dto.CategoryAddInput{
			TelegramID: update.Message.From.ID,
			Name:       name,
		})
		if err != nil {
			if err == domainerr.ErrCategoryAlreadyExists {
				return &transportdto.Result{
					ChatID: update.Message.Chat.ID,
					Text:   `❌ Категория "` + name + `" уже существует`,
				}, nil
			}
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Произошла ошибка. Попробуйте позже.",
			}, err
		}

		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text: "✅ Категория создана!\n\n" +
				"📂 Название: " + out.Name + "\n" +
				"🆔 ID: " + out.ID.String() + "\n\n" +
				"Используйте этот ID для удаления категории.",
		}, nil

	case "delete":
		if len(parts) < 2 {
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Укажите ID.\nПример: /category delete 00000000-0000-0000-0000-000000000000",
			}, nil
		}

		id, err := uuid.Parse(parts[1])
		if err != nil {
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Неверный формат ID.",
			}, nil
		}

		err = h.categoryService.Delete(ctx, dto.CategoryDeleteInput{
			TelegramID: update.Message.From.ID,
			CategoryID: id,
		})
		if err != nil {
			if err == domainerr.ErrCategoryNotFound {
				return &transportdto.Result{
					ChatID: update.Message.Chat.ID,
					Text:   "❌ Категория не найдена (или базовую категорию удалять нельзя).",
				}, nil
			}
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Произошла ошибка. Попробуйте позже.",
			}, err
		}

		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "✅ Категория удалена",
		}, nil

	default:
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Неизвестная подкоманда.\nИспользуйте:\n/category add <название>\n/category delete <id>",
		}, nil
	}
}
