package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	cfg := LoadConfig()

	if assert.NotEmpty(t, os.Getenv(cfg.Telegram.TokenEnv)) {
		fmt.Printf("os.Getenv(cfg.Telegram.TokenEnv): %v\n", os.Getenv(cfg.Telegram.TokenEnv))
	}
}
