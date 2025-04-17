package initial

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"video-realtime-ranking/config"
	"video-realtime-ranking/internal/app"
	"video-realtime-ranking/internal/dataaccess/redis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func Initial(cfg config.Config) {
	// signal notify
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	redis := redis.NewRedis(cfg)
	redisClient, err := redis.Connect()
	if err != nil {
		log.Fatal("Cannot connect to redis ", err)
	}
	defer redisClient.Close()

	handler := echo.New()
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	// waitGroup
	waitGroup, ctx := errgroup.WithContext(ctx)

	// custom host, port
	var opts app.Option
	if cfg.Server.Host != "" || cfg.Server.Port != "" {
		opts = app.Port(cfg.Server.Host, cfg.Server.Port)
	}
	app.NewServer(handler, waitGroup, ctx, opts)

	err = waitGroup.Wait()
	if err != nil {
		fmt.Println("error from wait group:", err.Error())
	}
}
