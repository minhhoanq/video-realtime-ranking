package initial

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/minhhoanq/video-realtime-ranking/common/logger"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/config"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/app"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/dataaccess/database"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/dataaccess/redis"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/handler/resful"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/routes"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/service"

	"go.uber.org/zap"
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

	logger.Setup(cfg.Server.Env, cfg.LogLevel)
	l := logger.NewWrapLogger(zap.DebugLevel, false)

	l.Info("start logger")
	_, err := database.New(cfg)
	if err != nil {
		log.Fatal("Cannot connect to database ", err)
	}
	// defer pg.Close()

	redisInstance := redis.NewRedis(cfg)
	redisClient, err := redisInstance.Connect()
	if err != nil {
		log.Fatal("Cannot connect to redis ", err)
	}
	defer redisClient.Close()

	rankingRedisDataAccessor := redis.NewRankingDataAccessor(redisClient)
	rankingService := service.NewRankingService(rankingRedisDataAccessor)
	rankingHandler := resful.NewHandler(rankingService)

	routes := routes.NewRouter(http.NewServeMux(), rankingHandler)

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
