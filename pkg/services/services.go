package services

import (
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"

	"google.golang.org/grpc"
)

type Services struct {
	GrpcServer *grpc.Server
	Config     *config.AppConfig
}

func NewServices(g *grpc.Server, c *config.AppConfig) *Services {
	return &Services{g, c}
}
