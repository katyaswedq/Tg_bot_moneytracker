package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"tgfin/internal/app"
	dbpg "tgfin/pkg/pg"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}

	cfg, err := dbpg.Load()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфиг: %v", err)
	}

	pool, err := dbpg.NewPool(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}
	defer pool.Close()

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}

	log.Printf("Бот авторизован как @%s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Получен сигнал остановки. Останавливаю получение обновлений...")
		bot.StopReceivingUpdates()
	}()

	app.Run(bot, pool, updates)

	log.Println("Бот остановлен корректно")
}
