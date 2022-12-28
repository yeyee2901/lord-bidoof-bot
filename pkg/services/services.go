package services

import (
	"context"
	"fmt"
	"strings"

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
