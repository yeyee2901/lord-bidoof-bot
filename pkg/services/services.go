package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/debug"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/telegram"
	telegrampb "github.com/yeyee2901/proto-lord-bidoof-bot/gen/go/telegram/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// for handling special characters, list them here. It will be replaced with
// empty character
const SPECIAL_CHARACTERS = ".{}[]!?"

type Services struct {
	BotAPI     *tgbotapi.BotAPI
	GrpcServer *grpc.Server
	DataSource *datasource.DataSource
}

func NewServices(g *grpc.Server, ds *datasource.DataSource, bot *tgbotapi.BotAPI) *Services {
	return &Services{bot, g, ds}
}

func (se *Services) InitServices() {
	telegrampb.RegisterTelegramServiceServer(se.GrpcServer, se)

	if se.DataSource.Config.Grpc.Mode != "production" {
		reflection.Register(se.GrpcServer)
		fmt.Println("REFLECTION ENABLED")
		debug.DebugStruct(se.GrpcServer.GetServiceInfo())
	}
}

func (se *Services) BotStatus(ctx context.Context, pbIn *telegrampb.BotStatusRequest) (*telegrampb.BotStatusResponse, error) {
	t := telegram.NewTelegramService(se.DataSource, se.BotAPI)

	if resp, err := t.GetBotStatus(ctx); err != nil {
		log.Error().Err(err).Msg("rpc.BotStatus.result")
		return nil, err
	} else {
		return &telegrampb.BotStatusResponse{
			Id:                      resp.Id,
			IsBot:                   resp.IsBot,
			FirstName:               resp.FirstName,
			Username:                resp.Username,
			CanJoinGroups:           resp.CanJoinGroups,
			CanReadAllGroupMessages: resp.CanReadAllGroupMessages,
			SupportsInlineQueries:   false,
		}, nil
	}
}

func (se *Services) SendMessage(ctx context.Context, pbIn *telegrampb.SendMessageRequest) (*telegrampb.SendMessageResponse, error) {
	t := telegram.NewTelegramService(se.DataSource, se.BotAPI)

	// sanitize input, replace special characters with ""
	var (
		msg       = pbIn.GetText()
		strReader = strings.NewReader(SPECIAL_CHARACTERS)
	)

	for i := 0; i < strReader.Len(); i++ {
		if b, err := strReader.ReadByte(); err == nil {
			msg = strings.ReplaceAll(msg, string(b), "")
		} else {
			log.Error().Err(err).Msg("rpc.SendMessage.specialCharacter")
			return nil, status.Error(codes.Internal, "Cannot parse special character list")
		}
	}

	// send the message
	if res, err := t.SendChat(ctx, pbIn.GetChatId(), msg, pbIn.GetUseMarkdown()); err != nil {
		log.Error().Err(err).Msg("rpc.SendMessage.result")
		return nil, err
	} else {
		return &telegrampb.SendMessageResponse{
			MessageId: res.MessageID,
			ChatId:    pbIn.GetChatId(),
			Recipient: res.Recipient,
		}, nil
	}
}

// Get list of private chats from database
func (se *Services) GetPrivateChat(ctx context.Context, pbIn *telegrampb.GetPrivateChatRequest) (*telegrampb.GetPrivateChatResponse, error) {
	filter := datasource.NewQueryFilter()

	// check query filters
	if chatId := pbIn.GetFilterChatId(); len(chatId) != 0 {
		filter["chat_id"] = chatId
	}

	if username := pbIn.GetFilterUsername(); len(username) != 0 {
		filter["username"] = username
	}

	// create task context
	thisCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	errChan := make(chan error)
	fatalErr := make(chan error)
	result := make(chan []datasource.PrivateChat)

	// dispatch job to goroutine
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fatalErr <- fmt.Errorf("%v", err)
			}
		}()

		var dbRes []datasource.PrivateChat
		var err error
		if len(filter) != 0 {
			dbRes, err = se.DataSource.GetPrivateChatWithQueryFilter(filter)
		} else {
			dbRes, err = se.DataSource.GetPrivateChatWithQueryFilter(nil)
		}

		if err != nil {
			errChan <- err
		} else {
			result <- dbRes
		}
	}()

	// poll for goroutine result
	for {
		select {
		// timeout or somehow this context got canceled
		case <-thisCtx.Done():
			if err := thisCtx.Err(); err == context.DeadlineExceeded {
				log.Error().Err(err).Msg("rpc.GetPrivateChat.timeout")
				return nil, status.Error(codes.DeadlineExceeded, "RPC timeout")
			} else {
				log.Error().Err(err).Msg("rpc.GetPrivateChat.canceled")
				return nil, status.Error(codes.Canceled, "RPC canceled by server")
			}

		// something happened when fetching from database
		case err := <-errChan:
			log.Error().Err(err).Msg("rpc.GetPrivateChat.database")
			return nil, status.Error(codes.Internal, "An error occured when querying to database")

		// fatal error happened
		case err := <-fatalErr:
			log.Error().Err(err).Msg("rpc.GetPrivateChat.FATAL")
			return nil, status.Error(codes.Internal, "Fatal internal server error")

		// successful case
		case res := <-result:
			log.Info().Interface("db_result", res).Msg("rpc.GetPrivateChat.result")

			// check if there's any result to avoid nil pointer deref
			if len(res) == 0 {
				return nil, status.Error(codes.NotFound, "No chats found.")
			}

			// iterate to assign values
			pbOut := &telegrampb.GetPrivateChatResponse{}
			for i := range res {
				chatData := &telegrampb.ChatData{
					ChatId:      res[i].ChatID,
					Username:    res[i].Username,
					DisplayName: res[i].Name,
					Bio:         res[i].Bio,
				}

				pbOut.Data = append(pbOut.Data, chatData)
			}
			pbOut.Count = uint64(len(res))

			return pbOut, nil
		}
	}
}
