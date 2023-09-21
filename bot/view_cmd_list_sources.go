package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sweetheart11/tgbot/botkit"
	"github.com/Sweetheart11/tgbot/lib/escmarkdn"
	"github.com/Sweetheart11/tgbot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

type SourceLister interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

func ViewCmdListSources(lister SourceLister) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := lister.Sources(ctx)
		if err != nil {
			return err
		}

		var (
			sourceInfos = lo.Map(sources, func(source model.Source, _ int) string {
				return formatSource(source)
			})
			msgText = fmt.Sprintf(
				"Available sources \\(total %d\\):\n\n%s",
				len(sources),
				strings.Join(sourceInfos, "\n\n"),
			)
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = "MarkdownV2"

		_, err = bot.Send(reply)

		return err
	}
}

func formatSource(source model.Source) string {
	return fmt.Sprintf(
		"🌐 *%s*\nID: '%d'\nURL: %s",
		escmarkdn.EscapeForMarkdown(source.Name),
		source.ID,
		escmarkdn.EscapeForMarkdown(source.FeedURL),
	)
}
