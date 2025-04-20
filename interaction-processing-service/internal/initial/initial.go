package initial

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/config"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/app"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/dataaccess/database"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/dataaccess/kafka/consumer"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/dataaccess/kafka/producer"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/dataaccess/redis"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/handler/consumers"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/handler/resful"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/routes"
	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/service"

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

	kafkaProducer, err := producer.NewProducer(cfg)
	if err != nil {
		log.Fatal("cannot connect to kafka producer ", err)
	}
	interactionCreateKafkaProducer := producer.NewInteractionProducer(kafkaProducer)

	kafkaConsumer, err := consumer.NewConsumer(cfg)
	if err != nil {
		log.Fatal("cannot connect to kafka consumer ", err)
	}

	interactionDataAccessor := database.NewInteractionDataAccessor(db)
	interactionService := service.NewInteractionService(interactionDataAccessor)
	interactionHandler := resful.NewHandler(interactionService, interactionCreateKafkaProducer)
	routes := routes.NewRouter(http.NewServeMux(), interactionHandler)

	// waitGroup
	waitGroup, ctx := errgroup.WithContext(ctx)

	// custom host, port
	var opts app.Option
	if cfg.Server.Host != "" || cfg.Server.Port != "" {
		opts = app.Port(cfg.Server.Host, cfg.Server.Port)
	}
	app.NewServer(routes.SetupRouter(), waitGroup, ctx, opts)

	// start consumers
	interactionProcessedMessageHandler := consumers.NewInteractionProcessedMessageHandler(interactionService)
	newInteractionKafkaConsumer := consumers.NewInteractionServiceKafkaConsumer(kafkaConsumer, interactionProcessedMessageHandler)

	go func(ctx context.Context) {
		err := newInteractionKafkaConsumer.Start(ctx)
		if err != nil {
			return
		}
	}(context.Background())

	err = waitGroup.Wait()
	if err != nil {
		fmt.Println("error from wait group:", err.Error())
	}
}
