package initial

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"video-realtime-ranking/interaction-processing-service/config"
	"video-realtime-ranking/interaction-processing-service/internal/app"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/database"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/redis"
	"video-realtime-ranking/interaction-processing-service/internal/handler/resful"
	"video-realtime-ranking/interaction-processing-service/internal/routes"
	"video-realtime-ranking/interaction-processing-service/internal/service"

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

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("Cannot connect to database ", err)
	}
	// defer mongod.Disconnect(context.Background())

	redis := redis.NewRedis(cfg)
	redisClient, err := redis.Connect()
	if err != nil {
		log.Fatal("Cannot connect to redis ", err)
	}
	defer redisClient.Close()

	interactionDataAccessor := database.NewInteractionDataAccessor(db)
	interactionService := service.NewInteractionService(interactionDataAccessor)
	interactionHandler := resful.NewHandler(interactionService)
	routes := routes.NewRouter(http.NewServeMux(), interactionHandler)

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
