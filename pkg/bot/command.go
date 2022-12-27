package bot

import (
	"context"
	"database/sql"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

// for now, bot can only be used in private chats
func (tg *TelegramBotService) StartCommand(ctx context.Context, msg *tgbotapi.Message, _ []string) {
	chat := msg.Chat

	// check if its private chat, if yes do nothing
	if !chat.IsPrivate() {
		tg.SendNormalChat(chat.ID, "Bidoof would like to apologize, but currently I cannot handle group chats for I am anti-social", "StartCommand.IsPrivate")
		return
	}

	// check if user chat id is already registered
	_, err := tg.GetPrivateChat(chat.ID)
	if err == nil {
		text := msg.From.FirstName + ", looks like you've already awaken Grand Lord Bidoof!"
		tg.SendNormalChat(chat.ID, text, "StartCommand.GetPrivateChat")
		return
	}

	switch {
	// user has not started the bot yet, so register them
	case err == sql.ErrNoRows:
		tg.savePrivateChat(chat)

		// inform user
		text := msg.From.FirstName + ", thank you for waking me. Bidoof bless you."
		tg.SendNormalChat(chat.ID, text, "StartCommand.savePrivateChat")

	// system error (db)
	case err != nil:
		panic(err)
	}
}

// Say something to another user via bidoof
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
