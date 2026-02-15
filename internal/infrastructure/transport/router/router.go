package router

import (
	"context"

	"tgfin/internal/infrastructure/transport/handlers"
	"tgfin/internal/infrastructure/transport/transportdto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	startHandler      *handlers.StartHandler
	helpHandler       *handlers.HelpHandler
	categoryHandler   *handlers.CategoryHandler
	categoriesHandler *handlers.CategoriesHandler
	addExpenseHandler *handlers.AddExpenseHandler
}

func NewRouter(startHandler *handlers.StartHandler, helpHandler *handlers.HelpHandler, categoryHandler *handlers.CategoryHandler, categoriesHandler *handlers.CategoriesHandler, addExpenseHandler *handlers.AddExpenseHandler) *Router {
	return &Router{
		startHandler:      startHandler,
		helpHandler:       helpHandler,
		categoryHandler:   categoryHandler,
		categoriesHandler: categoriesHandler,
		addExpenseHandler: addExpenseHandler,
	}
}

func (r *Router) Route(ctx context.Context, update tgbotapi.Update) (*transportdto.Result, error) {
	if update.Message == nil {
		return nil, nil
	}

	switch update.Message.Command() {
	case "start":
		return r.startHandler.Handle(ctx, update)
	case "help":
		return r.helpHandler.Handle(ctx, update)
	case "categories":
		return r.categoriesHandler.Handle(ctx, update)
	case "category":
		return r.categoryHandler.Handle(ctx, update)
	case "add":
		return r.addExpenseHandler.Handle(ctx, update)
	default:
		return &transportdto.Result{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Неизвестная команда. Используйте /help",
		}, nil
	}
}
