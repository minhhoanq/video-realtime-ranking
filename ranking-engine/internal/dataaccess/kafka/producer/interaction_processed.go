package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	TOPIC_NAME_INTERACTION_SERVICE_INTERACTION_PROCESSED = "interaction_service_interaction_processed"
)

type InteractionProcessed struct {
	VideoID         uuid.UUID
	UserID          uuid.UUID
	InteractionType string
}

type RankingProducer interface {
	Produce(ctx context.Context, message InteractionProcessed) error
}

type rankingProducer struct {
	producer Producer
}

func NewRankingProducer(producer Producer) RankingProducer {
	return &rankingProducer{
		producer: producer,
	}
}

func (rp *rankingProducer) Produce(ctx context.Context, message InteractionProcessed) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("failed to marshal message, err: ", err.Error())
		return err
	}

	if err = rp.producer.Produce(ctx, TOPIC_NAME_INTERACTION_SERVICE_INTERACTION_PROCESSED, messageBytes); err != nil {
		fmt.Println("failed to produce interaction create event, err: ", err.Error())
		return err
	}

	return nil
}
