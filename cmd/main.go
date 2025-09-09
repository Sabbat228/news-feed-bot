package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"news-feed-bot/internal/bot"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/config"
	fetcher2 "news-feed-bot/internal/fetcher"
	notifier2 "news-feed-bot/internal/notifier"
	"news-feed-bot/internal/storage"
	"news-feed-bot/internal/summary"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Println(err)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}
	defer db.Close()
	var (
		articleStorage = storage.NewArticlePostgresStorage(db)
		sourceStorage  = storage.NewSourcePostgresStorage(db)
		fetcher        = fetcher2.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		notifier = notifier2.New(
			articleStorage,
			summary.NewOpenAiSummarizer(config.Get().OpenAIKey, config.Get().OpenAIPrompt),
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStar())

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] fetcher error: %v", err)
			}

			log.Printf("fetcher stopped: %v", err)
		}
	}(ctx)
	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] notifier error: %v", err)
				return
			}
			log.Printf("notifier stopped: %v", err)
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] news bot: %v", err)
			return
		}
		log.Printf("news bot stopped: %v", err)
	}

}
