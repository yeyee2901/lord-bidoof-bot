package services

import (
	"context"

	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
	telegrampb "github.com/yeyee2901/proto-lord-bidoof-bot/gen/go/telegram/v1"

	"google.golang.org/grpc"
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
}

func (se *Services) BotStatus(ctx context.Context, pbIn *telegrampb.BotStatusRequest) (*telegrampb.BotStatusResponse, error) {
    
}
