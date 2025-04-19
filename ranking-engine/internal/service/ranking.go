package service

import (
	"context"
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
	_, err := i.rankingRedisDataAccessor.IncrementScore(ctx, redis.VideoLeaderBoardKey, arg.VideoID.String(), scores[arg.InteractionType])
	if err != nil {
		return err
	}

	message := producer.InteractionProcessed{
		VideoID:         arg.VideoID,
		UserID:          arg.UserID,
		InteractionType: arg.InteractionType,
	}
	return i.kafkaProducer.Produce(ctx, message)
}
