package consumers

import (
	"context"
	"encoding/json"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/kafka/consumer"
)

type OrderServiceKafkaConsumer interface {
	Start(ctx context.Context) error
}

type orderServiceKafkaConsumer struct {
	kafkaConsumer consumer.Consumer
}

func NewOrderServiceKafkaConsumer(kafkaConsumer consumer.Consumer) OrderServiceKafkaConsumer {
	return &orderServiceKafkaConsumer{
		kafkaConsumer: kafkaConsumer,
	}
}

func (o orderServiceKafkaConsumer) Start(ctx context.Context) error {
	o.kafkaConsumer.RegisterHandler(
		"TOPIC_NAME_PAYMENT_SERVICE_TRANSACTION_COMPLETED",
		func(ctx context.Context, topic string, message []byte) error {
			var payload any
			if err := json.Unmarshal(message, &payload); err != nil {
				// o.l.Error("failed to unmarshal message", zap.Error(err))
				return err
			}

			// o.paymentTransactionCompleted.Handle(ctx, payload)
			return nil
		},
	)

	return o.kafkaConsumer.Start(ctx)
}
