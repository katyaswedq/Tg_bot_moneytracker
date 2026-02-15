package app

import (
	"context"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	repoPG "tgfin/internal/infrastructure/repo/pg"
	"tgfin/internal/infrastructure/transport/handlers"
	"tgfin/internal/infrastructure/transport/router"
	"tgfin/internal/service"
	dbpg "tgfin/pkg/pg"
)

func Run(bot *tgbotapi.BotAPI, pool *pgxpool.Pool, updates tgbotapi.UpdatesChannel) {
	pgClient := &dbpg.Client{Pool: pool}

	userRepo := repoPG.NewUserRepo()
	categoryRepo := repoPG.NewCategoryRepo()
	expenseRepo := repoPG.NewExpenseRepo()

	startService := service.NewStartService(pgClient, userRepo, categoryRepo)
	startHandler := handlers.NewStartHandler(startService)

	helpHandler := handlers.NewHelpHandler()

	categoryService := service.NewCategoryService(pgClient, userRepo, categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	categoriesHandler := handlers.NewCategoriesHandler(categoryService)

	expenseService := service.NewExpenseService(pgClient, userRepo, categoryRepo, expenseRepo)
	addExpenseHandler := handlers.NewAddExpenseHandler(expenseService)

	r := router.NewRouter(startHandler, helpHandler, categoryHandler, categoriesHandler, addExpenseHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updatesIn := make(chan Task, 100)
	deliveries := make(chan Delivery, 100)

	var senderWG sync.WaitGroup
	senderWG.Add(1)
	go func() {
		defer senderWG.Done()
		for d := range deliveries {
			if d.Text == "" {
				continue
			}
			msg := tgbotapi.NewMessage(d.ChatID, d.Text)
			if _, err := bot.Send(msg); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
			}
		}
	}()

	const workers = 10
	var workersWG sync.WaitGroup
	workersWG.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer workersWG.Done()

			for task := range updatesIn {
				upd := task.Update
				if upd.Message == nil {
					continue
				}

				if !upd.Message.IsCommand() {
					continue
				}

				res, err := r.Route(ctx, upd)
				if err != nil {
					log.Println("Ошибка обработки команды:", err)
					continue
				}
				if res == nil {
					continue
				}

				deliveries <- Delivery{
					ChatID: res.ChatID,
					Text:   res.Text,
				}
			}
		}()
	}

	for update := range updates {
		updatesIn <- Task{Update: update}
	}

	close(updatesIn)
	workersWG.Wait()
	close(deliveries)
	senderWG.Wait()
}
