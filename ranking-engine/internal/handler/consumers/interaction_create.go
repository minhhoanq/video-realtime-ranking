package consumers

import (
	"context"
	"video-realtime-ranking/ranking-engine/internal/service"

	"github.com/google/uuid"
)

const (
	TOPIC_NAME_RANKING_SERVICE_INTERACTION_CREATE = "ranking_service_interaction_create"
)

type InteractionCreate struct {
	VideoID         uuid.UUID
	UserID          uuid.UUID
	InteractionType string
}

type InteractionCreateMessageHandler interface {
	Handle(ctx context.Context, message InteractionCreate) error
}

type interactionCreateMessageHandler struct {
	rankingService service.RankingEngineService
}

func NewInteractionCreateMessageHandler(rankingService service.RankingEngineService) InteractionCreateMessageHandler {
	return &interactionCreateMessageHandler{
		rankingService: rankingService,
	}
}

func (r *interactionCreateMessageHandler) Handle(ctx context.Context, message InteractionCreate) error {
	arg := &service.InteractionRequest{
		VideoID:         message.VideoID,
		UserID:          message.UserID,
		InteractionType: message.InteractionType,
	}

	return r.rankingService.RankingEngine(ctx, arg)
}
