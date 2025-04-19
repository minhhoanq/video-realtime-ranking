package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	TOPIC_NAME_INTERACTION_SERVICE_INTERACTION_CREATE = "interaction_service_interaction_create"
)

type InteractionCreate struct {
	VideoID         uuid.UUID
	UserID          uuid.UUID
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

	if err = ip.producer.Produce(ctx, TOPIC_NAME_INTERACTION_SERVICE_INTERACTION_CREATE, messageBytes); err != nil {
		fmt.Println("failed to produce interaction create event, err: ", err.Error())
		return err
	}

	return nil
}
