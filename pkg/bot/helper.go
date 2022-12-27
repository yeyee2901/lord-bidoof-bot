package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
)

// save incoming private chat
func (tg *TelegramBotService) savePrivateChat(chat *tgbotapi.Chat) {
	newChat := &datasource.PrivateChat{
		ChatID:   chat.ID,
		Username: chat.UserName,
		Name:     chat.FirstName + " " + chat.LastName,
		Bio:      chat.Bio,
	}

	if err := tg.InsertPrivateChatToDB(newChat); err != nil {
		panic(err)
	}
}
