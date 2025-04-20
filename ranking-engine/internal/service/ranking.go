package service

import (
	"context"
	"fmt"
	"video-realtime-ranking/ranking-engine/internal/dataaccess/database"
	"video-realtime-ranking/ranking-engine/internal/dataaccess/kafka/producer"
	"video-realtime-ranking/ranking-engine/internal/dataaccess/redis"

	"github.com/google/uuid"
)

type RankingEngineService interface {
	RankingEngine(ctx context.Context, arg *InteractionRequest) error
}

type rankingEngineService struct {
	interactionDataAccessor  database.InteractionDataAccessor
	rankingRedisDataAccessor redis.RankingDataAccessor
	kafkaProducer            producer.RankingProducer
}

func NewrankingEngineService(
	interactionDataAccessor database.InteractionDataAccessor,
	rankingRedisDataAccessor redis.RankingDataAccessor,
	kafkaProducer producer.RankingProducer,
) RankingEngineService {
	return &rankingEngineService{
		interactionDataAccessor:  interactionDataAccessor,
		rankingRedisDataAccessor: rankingRedisDataAccessor,
		kafkaProducer:            kafkaProducer,
	}
}

type InteractionRequest struct {
	VideoID         uuid.UUID
	UserID          uuid.UUID
	InteractionType string
}

func (i *rankingEngineService) RankingEngine(ctx context.Context, arg *InteractionRequest) error {
	scores := map[string]float64{
		"view":    1.0,
		"like":    1.0,
		"comment": 2.0,
		"share":   5.0,
	}

	score, ok := scores[arg.InteractionType]
	if !ok {
		return fmt.Errorf("invalid interaction type: %s", arg.InteractionType)
	}

	// increment score for global score and per user score
	scoreKeys := []string{
		redis.GlobalScore,
		fmt.Sprintf("%s:%s", redis.UserVideoScore, arg.UserID.String()),
	}

	for _, key := range scoreKeys {
		if _, err := i.rankingRedisDataAccessor.IncrementScore(ctx, key, arg.VideoID.String(), score); err != nil {
			return err
		}
	}

	message := producer.InteractionProcessed{
		VideoID:         arg.VideoID,
		UserID:          arg.UserID,
		InteractionType: arg.InteractionType,
	}
	return i.kafkaProducer.Produce(ctx, message)
}
