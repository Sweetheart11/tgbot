package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Sweetheart11/tgbot/botkit"
	"github.com/Sweetheart11/tgbot/lib/escmarkdn"
	"github.com/Sweetheart11/tgbot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceDeleter interface {
	Delete(ctx context.Context, id int64) error
	SourceByID(ctx context.Context, id int64) (*model.Source, error)
}

func ViewCmdDeleteSource(deleter SourceDeleter) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		id, err := strconv.ParseInt(update.Message.CommandArguments(), 10, 64)

		source, err := deleter.SourceByID(ctx, id)
		if err != nil {
			return err
		}

		if err := deleter.Delete(ctx, id); err != nil {
			return err
		}

		msgText := fmt.Sprintf(
			"Source deleted,\n%s\t:\t%s\n\n",
			source.Name,
			escmarkdn.EscapeForMarkdown(source.FeedURL),
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = "MarkdownV2"

		_, err = bot.Send(reply)

		return err
	}
}
