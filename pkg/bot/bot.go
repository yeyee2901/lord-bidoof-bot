package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command func(context.Context, []string) error

type TelegramBotService struct {
	*tgbotapi.BotAPI
	*datasource.DataSource
	Commands map[string]Command
}

func NewTelegramBotService(bot *tgbotapi.BotAPI, ds *datasource.DataSource) *TelegramBotService {
	tg := new(TelegramBotService)

	tg.DataSource = ds
	tg.BotAPI = bot
	tg.Commands = map[string]Command{
		"hello": tg.HelloCommand,
	}

	return tg
}

// handle update in separate goroutine so I can implement task timeouts
func (tg *TelegramBotService) HandleUpdate(event tgbotapi.Update) {
	updateCtx, cancel := context.WithTimeout(context.Background(), time.Duration(tg.Config.Telegram.Bot.Timeout)*time.Second)
	defer cancel()
	updateDone := make(chan struct{})

	// task goroutine
	go func() {
		defer func() {
			updateDone <- struct{}{}
		}()

		if event.Message.IsCommand() {
			if handler, exist := tg.Commands[event.Message.Command()]; exist {
				log.Info().Str("command", event.Message.Command()).Msg("handle.command")
				handler(updateCtx, strings.Split(event.Message.CommandArguments(), " "))
			} else {
				log.Warn().Str("command", event.Message.Command()).Msg("unknown.command")

				// inform user it was unknown command
				msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Unknown command")
				if _, err := tg.BotAPI.Send(msg); err != nil {
					log.Error().Err(err).Msg("send.error")
				}
			}
		}
	}()

	// wait for results
	for {
		select {
		case <-updateCtx.Done():
			if updateCtx.Err() == context.DeadlineExceeded {
				fmt.Println("timeout exceeded")
			} else {
				fmt.Println("Normal cancellation")
			}
			return

		case <-updateDone:
			fmt.Println("task done")
			return
		}
	}
}
