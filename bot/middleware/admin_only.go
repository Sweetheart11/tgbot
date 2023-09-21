package middleware

import (
	"context"

	"github.com/Sweetheart11/tgbot/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AdminOnly(channelID int64, next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: channelID,
			},
		})
		if err != nil {
			return err
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
				return next(ctx, bot, update)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You don't have the right to do this")); err != nil {
			return err
		}

		return nil
	}
}
