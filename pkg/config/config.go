package config

import (
	"os"
	"strings"

	"github.com/yeyee2901/lord-bidoof-bot/pkg/debug"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Grpc     grpcMeta     `yaml:"grpc"`
	Telegram telegramMeta `yaml:"telegram"`
	Redis    redisMeta    `yaml:"redis"`
	DB       databaseMeta `yaml:"db"`
}

type grpcMeta struct {
	Listener string `yaml:"listener"`
	Timeout  int    `yaml:"timeout"`
	Mode     string `yaml:"mode"`
	Logfile  string `yaml:"logfile"`
}

type telegramMeta struct {
	TokenEnv string  `yaml:"token_env"`
	Bot      botMeta `yaml:"bot"`
}

type botMeta struct {
	Logfile  string     `yaml:"logfile"`
	Timeout  int        `yaml:"timeout"`
	Messages botMessage `yaml:"messages"`
}

type botMessage struct {
	Panic          string `yaml:"panic"`
	UnknownCommand string `yaml:"unknown_command"`
	GroupChat      string `yaml:"group_chat"`
}

type redisMeta struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type databaseMeta struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Minpool  int    `yaml:"minpool"`
	Maxpool  int    `yaml:"maxpool"`
}

func LoadConfig() (config AppConfig) {
	b, err := os.ReadFile("setting/setting.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, &config); err != nil {
		panic(err)
	}

	// load token to environment
	b, err = os.ReadFile(".telegram-token")
	if err != nil {
		panic(err)
	}

	r := strings.NewReplacer("\n", "", "\r", "")
	token := r.Replace(string(b))
	if err = os.Setenv(config.Telegram.TokenEnv, token); err != nil {
		panic(err)
	}

	if config.Grpc.Mode != "production" {
		debug.DebugStruct(config)
	}

	return config
}
