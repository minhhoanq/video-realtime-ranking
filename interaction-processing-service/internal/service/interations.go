package service

import (
	"context"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/database"
)

type InteractionService interface {
	CreateInteraction(ctx context.Context, arg *database.SendInteractionRequest) (*database.SendInteractionResponse, error)
}

type interactionService struct {
	interactionDataAccessor database.InteractionDataAccessor
}

func NewInteractionService(
	interactionDataAccessor database.InteractionDataAccessor,
) InteractionService {
	return &interactionService{
		interactionDataAccessor: interactionDataAccessor,
	}
}

func (i *interactionService) CreateInteraction(ctx context.Context, arg *database.SendInteractionRequest) (*database.SendInteractionResponse, error) {
	interaction, err := i.interactionDataAccessor.CreateInteraction(ctx, arg)
	if err != nil {
		return nil, err
	}

	response := &database.SendInteractionResponse{
		UserID:  interaction.UserID,
		VideoID: interaction.UserID,
		ID:      interaction.ID,
	}

	return response, nil
}
