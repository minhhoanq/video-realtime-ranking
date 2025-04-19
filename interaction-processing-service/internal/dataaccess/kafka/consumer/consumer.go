package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"video-realtime-ranking/interaction-processing-service/config"

	"github.com/IBM/sarama"
)

type MessageHandlerFunc func(ctx context.Context, topic string, message []byte) error

type Consumer interface {
	RegisterHandler(topic string, handlerFunc MessageHandlerFunc)
	Start(ctx context.Context) error
}

type consumer struct {
	saramaConsumer        sarama.ConsumerGroup
	topicToHandlerFuncMap map[string]MessageHandlerFunc
}

func NewConsumer(cfg config.Config) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.ClientID = cfg.Kafka.ClientID
	config.Metadata.Full = true
	config.Consumer.Fetch.Min = 1024 * 1024
	config.Consumer.MaxWaitTime = 500 * time.Millisecond

	saramaConsumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.ClientID, config)

	if err != nil {
		return nil, fmt.Errorf("failed to create sarama consumer: %w", err)
	}

	return &consumer{
		saramaConsumer:        saramaConsumer,
		topicToHandlerFuncMap: make(map[string]MessageHandlerFunc),
	}, nil
}

func (c *consumer) RegisterHandler(topic string, handlerFunc MessageHandlerFunc) {
	c.topicToHandlerFuncMap[topic] = handlerFunc
}

func (c *consumer) Start(ctx context.Context) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for topic, handlerFunc := range c.topicToHandlerFuncMap {
		go func(topic string, handlerFunc MessageHandlerFunc) {
			if err := c.saramaConsumer.Consume(context.Background(), []string{topic}, newConsumerHandler(handlerFunc, signalChan)); err != nil {
				fmt.Println("faild to consumer message from queue, err: ", err.Error())
			}
		}(topic, handlerFunc)
	}

	return nil
}

type consumerHandler struct {
	handlerFunc MessageHandlerFunc
	signalChan  chan os.Signal
}

func newConsumerHandler(
	handlerFunc MessageHandlerFunc,
	signalChan chan os.Signal,
) *consumerHandler {
	return &consumerHandler{
		handlerFunc: handlerFunc,
		signalChan:  signalChan,
	}
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (h consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				session.Commit()
			}
			fmt.Printf("Order Service - Processing order: %s\n", string(message.Value))
			if err := h.handlerFunc(session.Context(), message.Topic, message.Value); err != nil {
				return err
			}
		case <-h.signalChan:
			session.Commit()
			return nil
		}
	}
}

func (h consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
