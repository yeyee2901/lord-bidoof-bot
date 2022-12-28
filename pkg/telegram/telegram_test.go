package telegram

import (
	"context"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/debug"
)

func TestGetBotStatus(t *testing.T) {
	cfg := config.LoadConfig()
	ds := datasource.NewDataSource(&cfg, nil, nil)
	bot, err := tgbotapi.NewBotAPI(os.Getenv(cfg.Telegram.TokenEnv))
	if err != nil {
		t.Fatal(err)
	}
	tg := NewTelegramService(ds, bot)

	resp, err := tg.GetBotStatus(context.Background())
	if assert.Nil(t, err) {
		debug.DebugStruct(resp)
	} else {
		t.Fatal(err)
	}
}
