package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TelegramService struct {
	BotAPI *tgbotapi.BotAPI
	*datasource.DataSource
}

func NewTelegramService(ds *datasource.DataSource, bot *tgbotapi.BotAPI) *TelegramService {
	return &TelegramService{bot, ds}
}

// get the bot status
func (t *TelegramService) GetBotStatus(ctx context.Context) (*RespGetMe, error) {
	// goroutine goes brrrr
	thisCtx, cancel := context.WithTimeout(ctx, time.Duration(t.Config.Telegram.Bot.Timeout)*time.Second)
	defer cancel()
	errChan := make(chan error)
	result := make(chan *RespGetMe)

	// task goroutine
	go func() {
		if user, err := t.BotAPI.GetMe(); err != nil {
			errChan <- err
		} else {
			result <- &RespGetMe{
				Id:                      uint64(user.ID),
				IsBot:                   user.IsBot,
				FirstName:               user.FirstName,
				LastName:                user.LastName,
				Username:                user.UserName,
				CanJoinGroups:           user.CanJoinGroups,
				CanReadAllGroupMessages: user.CanReadAllGroupMessages,
			}
		}
	}()

	for {
		select {

		// task timeout, or somehow canceled
		case <-thisCtx.Done():
			if err := thisCtx.Err(); err == context.DeadlineExceeded {
				return nil, status.Error(codes.DeadlineExceeded, "Timeout")
			} else {
				return nil, status.Error(codes.Canceled, err.Error())
			}

		// error from telegram
		case err := <-errChan:
			return nil, status.Error(codes.Aborted, err.Error())

		case res := <-result:
			return res, nil
		}
	}
}

// send chat to user with `chatId`
func (t *TelegramService) SendChat(ctx context.Context, chatId int64, message string, useMarkdown bool) (*RespSendMessage, error) {
	// goroutine goes brrrr
	chatCtx, cancel := context.WithTimeout(ctx, time.Duration(t.Config.Telegram.Bot.Timeout)*time.Second)
	defer cancel()
	errChan := make(chan error)
	result := make(chan *RespSendMessage)

	// send chat task
	go func() {
		toSend := tgbotapi.NewMessage(chatId, message)

		if useMarkdown {
			toSend.ParseMode = tgbotapi.ModeMarkdownV2
		}

		if m, err := t.BotAPI.Send(toSend); err != nil {
			errChan <- err
		} else {
			result <- &RespSendMessage{
				MessageID: int64(m.MessageID),
				Recipient: fmt.Sprintf("%s %s", m.Chat.FirstName, m.Chat.LastName),
			}
		}
	}()

	// poll goroutine
	for {
		select {
		// timeout, or context get canceled somehow
		case <-chatCtx.Done():
			if err := chatCtx.Err(); err == context.DeadlineExceeded {
				return nil, status.Error(codes.DeadlineExceeded, "Timeout")
			} else {
				return nil, status.Error(codes.Canceled, err.Error())
			}

		// case: error, something happened, the source must be from Telegram API
		case err := <-errChan:
			return nil, status.Error(codes.Aborted, err.Error())

		// case: sending chat success
		case resp := <-result:
			return resp, nil
		}
	}
}
