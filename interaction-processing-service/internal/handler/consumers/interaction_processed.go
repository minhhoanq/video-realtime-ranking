package consumers

import (
	"context"
	"fmt"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/database"
	"video-realtime-ranking/interaction-processing-service/internal/service"

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

type InteractionProcessedMessageHandler interface {
	Handle(ctx context.Context, message InteractionProcessed) error
}

type interactionProcessedMessageHandler struct {
	interactionService service.InteractionService
}

func NewInteractionProcessedMessageHandler(interactionService service.InteractionService) InteractionProcessedMessageHandler {
	return &interactionProcessedMessageHandler{
		interactionService: interactionService,
	}
}

func (i *interactionProcessedMessageHandler) Handle(ctx context.Context, message InteractionProcessed) error {
	fmt.Println("process create interaction")
	arg := &database.SendInteractionRequest{
		VideoID:         message.VideoID.String(),
		UserID:          message.UserID.String(),
		InteractionType: message.InteractionType,
	}

	_, err := i.interactionService.CreateInteraction(ctx, arg)
	return err
}
