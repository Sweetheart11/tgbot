package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sweetheart11/tgbot/bot"
	"github.com/Sweetheart11/tgbot/bot/middleware"
	"github.com/Sweetheart11/tgbot/botkit"
	"github.com/Sweetheart11/tgbot/config"
	"github.com/Sweetheart11/tgbot/fetcher"
	"github.com/Sweetheart11/tgbot/notifier"
	"github.com/Sweetheart11/tgbot/storage"
	"github.com/Sweetheart11/tgbot/summarizer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("failed to create bot api: %v", err)
		return
	}

	connStr, err := config.PostgresConnStr()
	if err != nil {
		log.Printf("failed to get db connection string: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Printf("failed to connect to db: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcher        = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		notifier = notifier.New(
			articleStorage,
			summarizer.NewOpenAISummarizer(config.Get().OpenAIKey, config.Get().OpenAIPrompt),
			botAPI,
			config.Get().NotificationInterval,
			config.Get().LookupTimeWindow,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.NewBot(botAPI)
	newsBot.RegisterCmdView(
		"start",
		middleware.AdminOnly(config.Get().TelegramChannelID, bot.ViewCmdStart()),
	)
	newsBot.RegisterCmdView(
		"addsource",
		middleware.AdminOnly(config.Get().TelegramChannelID, bot.ViweCmdAddSource(sourceStorage)),
	)
	newsBot.RegisterCmdView(
		"listsources",
		middleware.AdminOnly(config.Get().TelegramChannelID, bot.ViewCmdListSources(sourceStorage)),
	)

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start fetcher: %v", err)
				return
			}

			log.Println("fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start notifier: %v", err)
				return
			}

			log.Println("notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("failed to start bot: %v", err)
			return
		}

		log.Println("bot stopped")
	}
}
