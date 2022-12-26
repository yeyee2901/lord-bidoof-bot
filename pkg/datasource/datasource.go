package datasource

import "github.com/yeyee2901/lord-bidoof-bot/pkg/config"

type DataSource struct {
	Config *config.AppConfig
}

func NewDataSource(c *config.AppConfig) *DataSource {
	return &DataSource{c}
}
