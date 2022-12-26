package services

import (
	"context"
	"fmt"

	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/debug"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/telegram"
	telegrampb "github.com/yeyee2901/proto-lord-bidoof-bot/gen/go/telegram/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Services struct {
	GrpcServer *grpc.Server
	DataSource *datasource.DataSource
}

func NewServices(g *grpc.Server, ds *datasource.DataSource) *Services {
	return &Services{g, ds}
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
	t := telegram.NewTelegramService(se.DataSource)

	if botStatus, err := t.GetBotStatus(ctx); err != nil {
		// check which one causes the error
		switch err := err.(type) {
		case *telegram.ServerError:
			return nil, status.Error(codes.Internal, err.Error())

		case *telegram.TelegramError:
			return nil, status.Error(err.GrpcCode, err.Error())

		default:
			return nil, status.Error(codes.Unknown, "Unknown / unhandled error"+err.Error())
		}
	} else {
		// sukses
		return &telegrampb.BotStatusResponse{
			Id:                      botStatus.Result.Id,
			IsBot:                   botStatus.Result.IsBot,
			FirstName:               botStatus.Result.FirstName,
			Username:                botStatus.Result.Username,
			CanJoinGroups:           botStatus.Result.CanJoinGroups,
			CanReadAllGroupMessages: botStatus.Result.CanReadAllGroupMessages,
			SupportsInlineQueries:   false,
		}, nil
	}
}
