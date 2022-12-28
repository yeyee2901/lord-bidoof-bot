package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/yeyee2901/lord-bidoof-bot/pkg/bot"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// create context for this bot process (App level)
	appContext, appDone := context.WithCancel(context.Background())
	defer appDone()

	// init sub services
	cfg := config.LoadConfig()
	initLogger(&cfg)
	ds := datasource.NewDataSource(&cfg, initDB(&cfg), initRedis(&cfg))

	// init bot service
	botApi, err := tgbotapi.NewBotAPI(os.Getenv(cfg.Telegram.TokenEnv))
	if err != nil {
		panic(err)
	}
	botServer := bot.NewTelegramBotService(botApi, ds)

	// setup update channel for polling
	updateConfig := tgbotapi.NewUpdate(0)
	updateChan := botServer.BotAPI.GetUpdatesChan(updateConfig)
	defer botServer.BotAPI.StopReceivingUpdates()

	// channels for propagating datas
	quit := make(chan struct{})
	fatalError := make(chan error, 1) // make it unbuffered so it won't block

	// start listening for udpates
	go func() {
		log.Info().Msg("START")

		// NOTE: if the bot panics, I want it to exit gracefully
		// while also logging it to logfile for auditing
		defer func() {
			if err := recover(); err != nil {
				fatalError <- fmt.Errorf("%v", err)
				close(fatalError)

				quit <- struct{}{}
				close(quit)
			}
		}()

		for {
			select {
			case <-appContext.Done(): // normal exit
				return

			case newEvent := <-updateChan:
				log.Info().Interface("event", newEvent).Msg("event.new")
				botServer.HandleUpdate(newEvent)
			}
		}
	}()

	<-quit
	if err := <-fatalError; err != nil {
		fmt.Println(err)
		log.Error().Err(err).Msg("FATAL")
	}

	if err := appContext.Err(); err != nil {
		fmt.Println(err)
		log.Warn().Err(err).Msg("CANCELED")
	}

	log.Info().Msg("SHUTTING-DOWN")
}

func initLogger(cfg *config.AppConfig) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(&lumberjack.Logger{
		Filename:   cfg.Telegram.Bot.Logfile,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.With().Timestamp().Logger()
}

func initRedis(cfg *config.AppConfig) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	if err := r.Ping().Err(); err != nil {
		panic(err)
	}

	return r
}

func initDB(cfg *config.AppConfig) *sqlx.DB {
	dbConfig := mysql.Config{
		User:                 cfg.DB.User,
		Passwd:               cfg.DB.Password,
		Net:                  "tcp",
		Addr:                 cfg.DB.Host,
		DBName:               cfg.DB.Database,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}

	d := sqlx.MustConnect("mysql", dbConfig.FormatDSN())

	if err := d.Ping(); err != nil {
		panic(err)
	}

	return d
}
