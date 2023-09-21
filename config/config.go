package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken     string        `hcl:"telegram_bot_token"    env:"TELEGRAM_BOT_TOKEN"`
	TelegramChannelID    int64         `hcl:"telegram_channel_id"   env:"TELEGRAM_CHANNEL_ID"`
	DatabaseDSN          string        `hcl:"database_dsn"          env:"DATABASE_DSN"          default:"postgres://postgres:postgres@localhost:5432/news_feed_bot?sslmode=disable"`
	FetchInterval        time.Duration `hcl:"fetch_interval"        env:"FETCH_INTERVAL"        default:"10m"`
	NotificationInterval time.Duration `hcl:"notification_interval" env:"NOTIFICATION_INTERVAL" default:"1m"`
	FilterKeywords       []string      `hcl:"filter_keywords"       env:"FILTER_KEYWORDS"`
	OpenAIKey            string        `hcl:"openai_key"            env:"OPENAI_KEY"`
	OpenAIPrompt         string        `hcl:"openai_prompt"         env:"OPENAI_PROMPT"`
	OpenAIModel          string        `hcl:"openai_model"          env:"OPENAI_MODEL"          default:"gpt-3.5-turbo"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})
		if err := loader.Load(); err != nil {
			log.Printf("failed to load config: %v", err)
		}
	})

	return cfg
}

func PostgresConnStr() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v\n", err)
	}

	host := os.Getenv("POSTGRESDB_HOST")
	port := os.Getenv("POSTGRESDB_PORT")
	password := os.Getenv("POSTGRESDB_PASSWORD")
	dbname := os.Getenv("POSTGRESDB_NAME")
	user := os.Getenv("POSTGRESDB_USER")
	sslmode := os.Getenv("POSTGRESDB_SSLMODE")
	//
	// connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
	// 	user, password, host, port, dbname)
	// connStr := "hostuser=user dbname=mydb password=pass sslmode=disable"
	connStr := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)

	return connStr, nil
}
