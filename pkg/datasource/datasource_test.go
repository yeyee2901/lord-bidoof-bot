package datasource

import (
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/debug"

	"github.com/jmoiron/sqlx"
)

func TestGetPrivateChatWithQueryFilter(t *testing.T) {
	cfg := config.LoadConfig()
	db := initDB(&cfg)
	ds := NewDataSource(&cfg, db, nil)

	testFilter := []struct {
		Name   string
		Filter QueryFilter
	}{
		{
			Name: "filter_chat_id",
			Filter: QueryFilter{
				"chat_id": "1234",
			},
		},
		{
			Name: "filter_username",
			Filter: QueryFilter{
				"username": "gabriel_s",
			},
		},
		{
			Name:   "filter_none",
			Filter: nil,
		},
	}

	for _, test := range testFilter {
		t.Run(test.Name, func(t *testing.T) {
			if res, err := ds.GetPrivateChatWithQueryFilter(test.Filter); err != nil {
				t.Fatal(err)
			} else {
				debug.DebugStruct(res)
			}
		})
	}
}

func initDB(cfg *config.AppConfig) *sqlx.DB {
	dbConfig := &mysql.Config{
		User:                 cfg.DB.User,
		Passwd:               cfg.DB.Password,
		Net:                  "tcp",
		Addr:                 cfg.DB.Host,
		DBName:               cfg.DB.Database,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}
	return sqlx.MustConnect("mysql", dbConfig.FormatDSN())
}
