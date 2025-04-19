package producer

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	TOPIC_NAME_RANKING_SERVICE_INTERACTION_CREATE = "ranking_service_interaction_create"
)

type InteractionCreate struct {
	VideoID         string
	UserID          string
	InteractionType string
}

type InteractionCreateProducer interface {
	Produce(ctx context.Context, message InteractionCreate) error
}

type interactionCreateProducer struct {
	producer Producer
}

func NewInteractionProducer(producer Producer) InteractionCreateProducer {
	return &interactionCreateProducer{
		producer: producer,
	}
}

func (ip *interactionCreateProducer) Produce(ctx context.Context, message InteractionCreate) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("failed to marshal message, err: ", err.Error())
		return err
	}

	if err = ip.producer.Produce(ctx, TOPIC_NAME_RANKING_SERVICE_INTERACTION_CREATE, messageBytes); err != nil {
		fmt.Println("failed to produce interaction create event, err: ", err.Error())
		return err
	}

	return nil
}
