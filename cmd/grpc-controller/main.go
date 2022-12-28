package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/services"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"
)

type App struct {
	Config     *config.AppConfig
	DB         *sqlx.DB
	Redis      *redis.Client
	GrpcServer *grpc.Server
}

func main() {
	app := InitApp()
	defer app.Cleanup()

	if err := app.Run(); err != nil {
		log.Info().AnErr("error", err).Msg("EXIT")
	}
}

func InitApp() *App {
	app := new(App)

	// INIT: config file
	cfg := config.LoadConfig()
	app.Config = &cfg

	// INIT: logger
	app.InitLogger()

	// INIT: db
	app.InitDB()

	// INIT: redis
	app.InitRedis()

	// INIT: gRPC server
	app.InitGrpc()

	return app
}

func (app *App) Run() error {
	log.Info().Msg("START")

	// channel for propagating handling error & OS interrupt
	errChan := make(chan error)
	fatalError := make(chan error, 1)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// create listener for grpc server to attach
	lst, err := net.Listen("tcp", app.Config.Grpc.Listener)
	if err != nil {
		return err
	}

	// run grpc server
	go func() {
		// make sure we stop the server to free the resource
		defer func() {
			if err := recover(); err != nil {
				fatalError <- fmt.Errorf("%v", err)
			}
			app.GrpcServer.GracefulStop()
			lst.Close()
		}()

		fmt.Println("Server listening at", app.Config.Grpc.Listener)
		errChan <- app.GrpcServer.Serve(lst)
	}()

	for {
		select {
		case <-sigChan:
			return fmt.Errorf("Server interrupted")

		case err := <-errChan:
			return err

		case err := <-fatalError:
			log.Error().Err(err).Msg("FATAL")
			return err
		}
	}
}

func (app *App) InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(&lumberjack.Logger{
		Filename:   app.Config.Grpc.Logfile,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.With().Timestamp().Logger()
}

func (app *App) InitDB() {
	dsn := mysql.Config{
		User:                 app.Config.DB.User,
		Passwd:               app.Config.DB.Password,
		Net:                  "tcp",
		Addr:                 app.Config.DB.Host,
		DBName:               app.Config.DB.Database,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}

	app.DB = sqlx.MustConnect("mysql", dsn.FormatDSN())
	if err := app.DB.Ping(); err != nil {
		panic(err)
	}
}

func (app *App) InitRedis() {
	app.Redis = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%s", app.Config.Redis.Host, app.Config.Redis.Port),
	})

	if err := app.Redis.Ping().Err(); err != nil {
		panic(err)
	}
}

func (app *App) InitGrpc() {
	app.GrpcServer = grpc.NewServer()
	ds := datasource.NewDataSource(app.Config, app.DB, app.Redis)

	// get bot token from environment
	token := os.Getenv(app.Config.Telegram.TokenEnv)
	if len(token) == 0 {
		panic("Empty bot token in environment variable")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	s := services.NewServices(app.GrpcServer, ds, bot)
	s.InitServices()
}

func (app *App) Cleanup() {
	if err := app.DB.Close(); err != nil {
		fmt.Println(err)
	}

	if err := app.DB.Close(); err != nil {
		fmt.Println(err)
	}

	if err := app.DB.Close(); err != nil {
		fmt.Println(err)
	}

}
