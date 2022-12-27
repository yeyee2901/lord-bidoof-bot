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

type Command func(context.Context, *tgbotapi.Message, []string)

type TelegramBotService struct {
	*tgbotapi.BotAPI
	*datasource.DataSource

	Commands map[string]Command
}

func NewTelegramBotService(bot *tgbotapi.BotAPI, ds *datasource.DataSource) *TelegramBotService {
	tg := new(TelegramBotService)

	tg.DataSource = ds
	tg.BotAPI = bot
	tg.RegisterCommands()

	return tg
}

// handle update in separate goroutine so I can implement task timeouts
func (tg *TelegramBotService) HandleUpdate(event tgbotapi.Update) {
	updateCtx, cancel := context.WithTimeout(context.Background(), time.Duration(tg.Config.Telegram.Bot.Timeout)*time.Second)
	defer cancel()
	updateDone := make(chan struct{})
	commandPanic := make(chan error)

	// task goroutine
	go func() {
		defer func() {
			// same as app level recovery, this handle command panics
			// and reports it to sender & logfile
			if err := recover(); err != nil {
				commandPanic <- fmt.Errorf("%v", err)
			} else {
				// normal flow
				updateDone <- struct{}{}
			}
		}()

		// check if its a command, otherwise do nothing
		if event.Message.IsCommand() {
			tg.handleCommand(updateCtx, event.Message)
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
			return

		case err := <-commandPanic:
			tg.handlePanic(err, event.Message)
			return
		}
	}
}

func (tg *TelegramBotService) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	// check if command exists
	handler, exist := tg.Commands[msg.Command()]
	if !exist {
		log.Warn().Str("command", msg.Command()).Msg("command.error")

		// inform user it was unknown command
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Unknown command")
		if _, err := tg.BotAPI.Send(msg); err != nil {
			log.Error().Err(err).Interface("message", msg).Msg("send.error")
		}

		return
	}

	// handle command normally, internal command error should not be informed
	// to user
	log.Info().Str("command", msg.Command()).Msg("command.handle")
	handler(ctx, msg, strings.Split(msg.CommandArguments(), " "))
}

func (tg *TelegramBotService) handlePanic(err error, msg *tgbotapi.Message) {
	log.Error().Err(err).Interface("message", msg).Msg("command.panic")

	text := fmt.Sprintf("I'm sorry, but Bidoof currently cannot process that, %s %s :(", msg.From.FirstName, msg.From.LastName)
	toUser := tgbotapi.NewMessage(msg.Chat.ID, text)

	if _, err := tg.BotAPI.Send(toUser); err != nil {
		log.Error().Err(err).Interface("message", msg).Msg("send.error")
	}
}
