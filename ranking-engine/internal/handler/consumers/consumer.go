package consumers

import (
	"context"
	"encoding/json"

	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/internal/dataaccess/kafka/consumer"
)

type RankingServiceKafkaConsumer interface {
	Start(ctx context.Context) error
}

type rankingServiceKafkaConsumer struct {
	kafkaConsumer     consumer.Consumer
	interactionCreate InteractionCreateMessageHandler
}

func NewRankingServiceKafkaConsumer(kafkaConsumer consumer.Consumer,
	interactionCreate InteractionCreateMessageHandler) RankingServiceKafkaConsumer {
	return &rankingServiceKafkaConsumer{
		kafkaConsumer:     kafkaConsumer,
		interactionCreate: interactionCreate,
	}
}

func (r *rankingServiceKafkaConsumer) Start(ctx context.Context) error {
	r.kafkaConsumer.RegisterHandler(
		TOPIC_NAME_RANKING_SERVICE_INTERACTION_CREATE,
		func(ctx context.Context, topic string, message []byte) error {
			var payload InteractionCreate
			if err := json.Unmarshal(message, &payload); err != nil {
				// o.l.Error("failed to unmarshal message", zap.Error(err))
				return err
			}

			r.interactionCreate.Handle(ctx, payload)
			return nil
		},
	)

	return r.kafkaConsumer.Start(ctx)
}
