package main

import (
	"fmt"
	"net"
	"time"

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
		panic(err)
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

func (app *App) Run() (err error) {
	if lst, err := net.Listen("tcp", app.Config.Grpc.Listener); err == nil {
		fmt.Println("Server listening at", app.Config.Grpc.Listener)
		err = app.GrpcServer.Serve(lst)
	}

	return err
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
	s := services.NewServices(app.GrpcServer, ds)
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
