package handlers

import (
	"context"
	"strconv"
	"strings"

	domainerr "tgfin/internal/domain/error"
	"tgfin/internal/infrastructure/transport/transportdto"
	"tgfin/internal/service/contract"
	"tgfin/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddExpenseHandler struct {
	expenseService contract.ExpenseService
}

func NewAddExpenseHandler(expenseService contract.ExpenseService) *AddExpenseHandler {
	return &AddExpenseHandler{
		expenseService: expenseService,
	}
}

func (h *AddExpenseHandler) Handle(ctx context.Context, update tgbotapi.Update) (*transportdto.Result, error) {
	if update.Message == nil || update.Message.From == nil {
		return nil, nil
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Использование:\n/add <сумма> <категория> <описание>",
		}, nil
	}

	parts := strings.Fields(args)
	if len(parts) < 2 {
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Использование:\n/add <сумма> <категория> <описание>",
		}, nil
	}

	amount, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Ошибка: сумма должна быть числом",
		}, nil
	}

	category := parts[1]

	var description *string
	if len(parts) > 2 {
		desc := strings.TrimSpace(strings.Join(parts[2:], " "))
		if desc != "" {
			description = &desc
		}
	}

	out, err := h.expenseService.Add(ctx, dto.ExpenseAddInput{
		TelegramID:   update.Message.From.ID,
		Amount:       amount,
		CategoryName: category,
		Description:  description,
	})
	if err != nil {
		switch err {
		case domainerr.ErrInvalidAmount:
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Сумма должна быть положительным числом",
			}, nil
		case domainerr.ErrCategoryNotFound:
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   `❌ Категория "` + category + `" не найдена`,
			}, nil
		default:
			return &transportdto.Result{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Произошла ошибка. Попробуйте позже.",
			}, err
		}
	}

	text := "✅ Расход добавлен!\n\n" +
		"💰 Сумма: " + strconv.FormatInt(out.Amount, 10) + "\n" +
		"📂 Категория: " + out.Category

	if out.Description != nil {
		text += "\n📝 Описание: " + *out.Description
	}

	return &transportdto.Result{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	}, nil
}
