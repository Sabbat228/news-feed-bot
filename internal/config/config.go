package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	TelegramBotToken     string        `hcl:"telegram_bot_token" env:"TELEGRAM_BOT_TOKEN" required:"true"`
	TelegramChannelID    int64         `hcl:"telegram_channel_id" env:"TELEGRAM_CHANNEL_ID" required:"true"`
	DatabaseDSN          string        `hcl:"database_dsn" env:"DATABASE_DSN" default:"jdbc:postgresql://localhost:5433/news_feed_bot"`
	FetchInterval        time.Duration `hcl:"fetch_interval" env:"FETCH_INTERVAL" default:"10m"`
	NotificationInterval time.Duration `hcl:"notification_interval" env:"NOIFICATION_INTERVAL" default:"1m"`
	FilterKeywords       []string      `hcl:"filter_keywords" env:"FILTER_KEYS"`
	OpenAIKey            string        `hcl:"open_ai_key" env:"OPEN_AI_KEY"`
	OpenAIPrompt         string        `hcl:"open_ai_prompt" env:"OPEN_AI_PROMPT"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		// Добавьте отладочную информацию
		wd, _ := os.Getwd()
		log.Printf("Текущая рабочая директория: %s", wd)

		// Проверьте существование файлов конфигурации
		for _, file := range []string{"./config.hcl", "./config.local.hcl"} {
			if _, err := os.Stat(file); err == nil {
				log.Printf("Файл найден: %s", file)
			} else {
				log.Printf("Файл не найден: %s", file)
			}
		}

		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] Failed to load config file: %v", err)
			return
		}

		// Отладочный вывод загруженных значений
		log.Printf("Загруженный TelegramBotToken: '%s'", cfg.TelegramBotToken)
		log.Printf("Загруженный TelegramChannelID: %d", cfg.TelegramChannelID)
	})
	return cfg
}
