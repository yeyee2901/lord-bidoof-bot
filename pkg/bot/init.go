package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg *TelegramBotService) InitBot() {
	tg.RegisterCommands()
}

func (tg *TelegramBotService) RegisterCommands() {
	tg.Commands = map[string]Command{
		"hello": tg.UnimplementedCommand,
		"start": tg.StartCommand,
		"stop":  tg.StopCommand,
	}

	cmdRegister := []tgbotapi.BotCommand{
		{
			Command:     "hello",
			Description: "Say something",
		},
		{
			Command:     "start",
			Description: "Start the bot",
		},
		{
			Command:     "stop",
			Description: "Stop bot interaction for this user",
		},
	}

	setBotCmd := tgbotapi.NewSetMyCommandsWithScope(tgbotapi.BotCommandScope{Type: "all_private_chats"}, cmdRegister...)
	if tgResp, err := tg.BotAPI.Request(setBotCmd); err != nil {
		panic(err)
	} else {
		if !tgResp.Ok {
			panic("Failed to register bot commands: " + tgResp.Description)
		}
	}
}
