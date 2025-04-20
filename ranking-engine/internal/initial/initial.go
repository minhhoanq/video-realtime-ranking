package initial

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/config"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/app"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/dataaccess/database"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/dataaccess/kafka/consumer"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/dataaccess/kafka/producer"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/dataaccess/redis"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/handler/consumers"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/routes"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/service"
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

	redisInstance := redis.NewRedis(cfg)
	redisClient, err := redisInstance.Connect()
	if err != nil {
		log.Fatal("Cannot connect to redis ", err)
	}
	defer redisClient.Close()

	rankingRepository := redis.NewRankingDataAccessor(redisClient)

	kafkaProducer, err := producer.NewProducer(cfg)
	if err != nil {
		log.Fatal("Cannot new producer kafka ", err)
	}
	rankingKafkaProducer := producer.NewRankingProducer(kafkaProducer)

	kafkaConsumer, err := consumer.NewConsumer(cfg)
	if err != nil {
		log.Fatal("Cannot new consumer kafka ", err)
	}

	interactionDataAccessor := database.NewInteractionDataAccessor(db)
	rankingService := service.NewrankingEngineService(interactionDataAccessor, rankingRepository, rankingKafkaProducer)
	routes := routes.NewRouter(http.NewServeMux(), rankingService)

	// waitGroup
	waitGroup, ctx := errgroup.WithContext(ctx)

	// custom host, port
	var opts app.Option
	if cfg.Server.Host != "" || cfg.Server.Port != "" {
		opts = app.Port(cfg.Server.Host, cfg.Server.Port)
	}
	app.NewServer(routes.SetupRouter(), waitGroup, ctx, opts)

	// start consumer threads
	interactionCreateMessageHandler := consumers.NewInteractionCreateMessageHandler(rankingService)
	newRankingServiceKafkaConsumer := consumers.NewRankingServiceKafkaConsumer(kafkaConsumer, interactionCreateMessageHandler)

	go func(ctx context.Context) {
		err := newRankingServiceKafkaConsumer.Start(ctx)
		if err != nil {
			return
		}
	}(context.Background())

	err = waitGroup.Wait()
	if err != nil {
		fmt.Println("error from wait group:", err.Error())
	}
}
