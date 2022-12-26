package telegram

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/debug"
)

func TestGetBotStatus(t *testing.T) {
	cfg := config.LoadConfig()
	ds := datasource.NewDataSource(&cfg, nil, nil)
	tg := NewTelegramService(ds)

	resp, err := tg.GetBotStatus(context.Background())
	if assert.Nil(t, err) {
		debug.DebugStruct(resp)
	} else {
		t.Fatal(err)
	}
}
