package database

import (
	"context"
	"fmt"
	"time"

	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/pkg/constants"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type InteractionDataAccessor interface {
	CreateInteraction(ctx context.Context, arg *SendInteractionRequest) (*SendInteractionResponse, error)
}

type interactionDataAccessor struct {
	database Database
}

func NewInteractionDataAccessor(database Database) InteractionDataAccessor {
	return &interactionDataAccessor{
		database: database,
	}
}

func (i *interactionDataAccessor) CreateInteraction(ctx context.Context, arg *SendInteractionRequest) (*SendInteractionResponse, error) {
	collection := i.database.returnCollectionPointer(constants.RANKING_COLLECTION)

	result, err := collection.InsertOne(ctx,
		Interaction{
			UserID:          arg.UserID,
			VideoID:         arg.VideoID,
			InteractionType: arg.InteractionType,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert notify to database, err: %s", err)
	}

	// convert insertedID type
	insertedID := result.InsertedID.(bson.ObjectID).Hex()

	response := &SendInteractionResponse{
		ID:      insertedID,
		UserID:  arg.UserID,
		VideoID: arg.VideoID,
	}

	return response, nil
}
