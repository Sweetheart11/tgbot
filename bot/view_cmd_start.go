package bot

import (
	"context"

	"github.com/Sweetheart11/tgbot/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Hello, World!")); err != nil {
			return err
		}

		return nil
	}
}
