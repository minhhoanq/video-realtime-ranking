package consumers

import (
	"context"
	"encoding/json"

	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/internal/dataaccess/kafka/consumer"
)

type InteractionServiceKafkaConsumer interface {
	Start(ctx context.Context) error
}

type interactionServiceKafkaConsumer struct {
	kafkaConsumer                      consumer.Consumer
	interactionProcessedMessageHandler InteractionProcessedMessageHandler
}

func NewInteractionServiceKafkaConsumer(kafkaConsumer consumer.Consumer,
	interactionProcessedMessageHandler InteractionProcessedMessageHandler) InteractionServiceKafkaConsumer {
	return &interactionServiceKafkaConsumer{
		kafkaConsumer:                      kafkaConsumer,
		interactionProcessedMessageHandler: interactionProcessedMessageHandler,
	}
}

func (i *interactionServiceKafkaConsumer) Start(ctx context.Context) error {
	i.kafkaConsumer.RegisterHandler(
		TOPIC_NAME_INTERACTION_SERVICE_INTERACTION_PROCESSED,
		func(ctx context.Context, topic string, message []byte) error {
			var payload InteractionProcessed
			if err := json.Unmarshal(message, &payload); err != nil {
				// o.l.Error("failed to unmarshal message", zap.Error(err))
				return err
			}

			return i.interactionProcessedMessageHandler.Handle(ctx, payload)
		},
	)

	return i.kafkaConsumer.Start(ctx)
}
