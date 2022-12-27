package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (tg *TelegramBotService) SendNormalChat(chatId int64, text, logSubject string) {
	msg := tgbotapi.NewMessage(chatId, text)
	if _, err := tg.BotAPI.Send(msg); err != nil {
		log.Error().Err(err).Interface("message", msg).Msg("send.error-" + logSubject)
	}
}
