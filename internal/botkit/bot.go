package botkit

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"runtime/debug"
	"time"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	cmdView map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}
	b.cmdView[cmd] = view
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if update.Message == nil || !update.Message.IsCommand() {
		return nil
	}

	cmd := update.Message.Command()
	view, ok := b.cmdView[cmd]
	if !ok {
		return nil
	}

	if err := view(ctx, b.api, update); err != nil {
		log.Printf("[ERROR] failed to handle command %s: %v", cmd, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при обработке команды.")
		if _, sendErr := b.api.Send(msg); sendErr != nil {
			log.Printf("[ERROR] failed to send error message: %v", sendErr)
		}
	}
	return nil
}
