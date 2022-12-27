package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/debug"
)

func (tg *TelegramBotService) RegisterCommands() {
	tg.Commands = map[string]Command{
		"hello": tg.HelloCommand,
	}

	cmdRegister := []tgbotapi.BotCommand{
		{
			Command:     "hello",
			Description: "say something",
		},
	}

	setBotCmd := tgbotapi.NewSetMyCommandsWithScope(tgbotapi.BotCommandScope{Type: "all_private_chats"}, cmdRegister...)
	if tgResp, err := tg.BotAPI.Request(setBotCmd); err != nil {
		panic(err)
	} else {
		debug.DebugStruct(tgResp)
		if !tgResp.Ok {
			panic("Failed to register bot commands: " + tgResp.Description)
		}
	}
}

func (tg *TelegramBotService) HelloCommand(ctx context.Context, msg *tgbotapi.Message, args []string) {
	// validate hello command
	if len(args) != 2 {
		tg.showUsage(msg.Chat.ID, HELLO_USAGE, false)
		return
	}

	baseStr := `
Hello %to% \! Bidoof wants to say: 

"%msg%"

That's all Bidoof have to say, sir\.
    `

	text := strings.Replace(baseStr, "%to%", args[0], 1)
	text = strings.Replace(text, "%msg%", args[1], 1)

	toUser := tgbotapi.NewMessage(msg.Chat.ID, text)
	toUser.ParseMode = tgbotapi.ModeMarkdownV2
	if _, err := tg.BotAPI.Send(toUser); err != nil {
		log.Error().Err(err).Msg("send.error.HelloCommand")
	}
}

func (tg *TelegramBotService) showUsage(chatId int64, usage string, useMarkdown bool) {
	msg := tgbotapi.NewMessage(chatId, usage)

	if useMarkdown {
		msg.ParseMode = tgbotapi.ModeMarkdownV2
	}

	if _, err := tg.BotAPI.Send(msg); err != nil {
		log.Error().Err(err).Msg("send.error.HelloCommand")
	}

	return
}
