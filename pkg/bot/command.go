package bot

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg *TelegramBotService) UnimplementedCommand(ctx context.Context, msg *tgbotapi.Message, _ []string) {
	panic("Unimplemented")
}

func (tg *TelegramBotService) StartCommand(ctx context.Context, msg *tgbotapi.Message, _ []string) {
	chat := msg.Chat

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

func (tg *TelegramBotService) StopCommand(ctx context.Context, msg *tgbotapi.Message, _ []string) {
	chat := msg.Chat

	switch _, err := tg.GetPrivateChat(chat.ID); {

	// no user found in DB, then do nothing
	case err == sql.ErrNoRows:
		text := "uh-oh, Who art thou? Zzzzz..."
		tg.SendNormalChat(chat.ID, text, "StopCommand.GetPrivateChat")
		return

	// system error (db)
	case err != nil:
		panic(err)

	// user found, then delete the chat record
	case err == nil:
		if err := tg.DeletePrivateChat(chat.ID); err != nil {
			panic(err)
		}

		text := fmt.Sprintf(`*Thank you for using me*\! If you need me, you can always /start me again or find me at t\.me/grandlordbidoof\_bot\. You can also safely delete this chat if you want\. Bidoof bless you\.`)
		tg.SendMarkdownChat(chat.ID, text, "StopCommand.DeletePrivateChat")
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

	tg.SendMarkdownChat(msg.Chat.ID, text, "HelloCommand")
}

func (tg *TelegramBotService) showUsage(chatId int64, usage string, useMarkdown bool) {
	if useMarkdown {
		tg.SendMarkdownChat(chatId, usage, "HelloCommand.showUsage")
	} else {
		tg.SendNormalChat(chatId, usage, "HelloCommand.showUsage")
	}

	return
}
