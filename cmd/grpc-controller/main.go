package main

import "github.com/yeyee2901/lord-bidoof-bot/pkg/config"

type App struct {
    Config *config.AppConfig
}

func main() {
    app := new(App)
    cfg := config.LoadConfig()
    
    app.Config = &cfg
}
