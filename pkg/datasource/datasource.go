package datasource

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"
)

type DataSource struct {
	Config *config.AppConfig
	DB     *sqlx.DB
	Redis  *redis.Client
}

func NewDataSource(c *config.AppConfig, db *sqlx.DB, r *redis.Client) *DataSource {
	return &DataSource{c, db, r}
}
