package botkit

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func NewBot(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}

	b.cmdViews[cmd] = view
}

func (b *Bot) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v\n%s", r, string(debug.Stack()))
		}
	}()

	var view ViewFunc

	if update.Message == nil {
		return
	}

	if !update.Message.IsCommand() {
		return
	}

	cmd := update.Message.Command()

	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		return
	}

	view = cmdView

	if err := view(ctx, b.api, update); err != nil {
		log.Printf("failed to handle update: %v\n", err)

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "internal error")); err != nil {
			log.Printf("failed to send error message: %v\n", err)
		}
	}
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			b.HandleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
