package initial

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"video-realtime-ranking/config"
	"video-realtime-ranking/internal/app"
	"video-realtime-ranking/internal/dataaccess/database"
	"video-realtime-ranking/internal/dataaccess/redis"
	"video-realtime-ranking/internal/routes"

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

	_, err := database.New(cfg)
	if err != nil {
		log.Fatal("Cannot connect to database ", err)
	}

	redis := redis.NewRedis(cfg)
	redisClient, err := redis.Connect()
	if err != nil {
		log.Fatal("Cannot connect to redis ", err)
	}
	defer redisClient.Close()

	routes := routes.NewRouter(http.NewServeMux())

	// waitGroup
	waitGroup, ctx := errgroup.WithContext(ctx)

	// custom host, port
	var opts app.Option
	if cfg.Server.Host != "" || cfg.Server.Port != "" {
		opts = app.Port(cfg.Server.Host, cfg.Server.Port)
	}
	app.NewServer(routes.SetupRouter(), waitGroup, ctx, opts)

	err = waitGroup.Wait()
	if err != nil {
		fmt.Println("error from wait group:", err.Error())
	}
}
