package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/minhhoanq/video-realtime-ranking/ranking-engine/config"
)

type Producer interface {
	Produce(ctx context.Context, topic string, message []byte) error
}

type producer struct {
	saramaSyncProducer sarama.SyncProducer
}

func NewProducer(cfg config.Config) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 1
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ClientID = cfg.Kafka.ClientID
	config.Metadata.Full = true

	saramaSyncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}

	return &producer{
		saramaSyncProducer: saramaSyncProducer,
	}, nil
}

func (p *producer) Produce(ctx context.Context, topic string, message []byte) error {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := p.saramaSyncProducer.SendMessage(&msg)
	if err != nil {
		return err
	}

	return nil
}
