package service

import (
	"context"
	"fmt"

	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/dataaccess/database"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/dataaccess/redis"
)

type RankingService interface {
	GetTopKVideoRanking(ctx context.Context, offset, limit int) ([]database.Score, error)
	GetTopKUserVideoRanking(ctx context.Context, user_id string, offset, limit int) ([]database.Score, error)
}

type rankingService struct {
	rankingRedisDataAccessor redis.RankingDataAccessor
}

func NewRankingService(
	rankingRedisDataAccessor redis.RankingDataAccessor,
) RankingService {
	return &rankingService{
		rankingRedisDataAccessor: rankingRedisDataAccessor,
	}
}

func (i *rankingService) GetTopKVideoRanking(ctx context.Context, offset, limit int) ([]database.Score, error) {
	start := int64(offset)
	stop := int64(offset + limit - 1)
	videos, err := i.rankingRedisDataAccessor.GetTopRanked(ctx, redis.GlobalScore, start, stop)
	if err != nil {
		return nil, err
	}

	videoScore := make([]database.Score, 0, len(videos))
	for _, video := range videos {
		videoScore = append(videoScore, database.Score{
			VideoID: video.Member.(string),
			Score:   video.Score,
		})
	}

	return videoScore, nil
}

func (i *rankingService) GetTopKUserVideoRanking(ctx context.Context, user_id string, offset, limit int) ([]database.Score, error) {
	start := int64(offset)
	stop := int64(offset + limit - 1)

	userVideokey := fmt.Sprintf("%s:%s", redis.UserVideoScore, user_id)
	videos, err := i.rankingRedisDataAccessor.GetTopRanked(ctx, userVideokey, start, stop)
	if err != nil {
		return nil, err
	}

	videoScore := make([]database.Score, 0, len(videos))
	for _, video := range videos {
		videoScore = append(videoScore, database.Score{
			VideoID: video.Member.(string),
			Score:   video.Score,
		})
	}

	return videoScore, nil
}
