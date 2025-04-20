package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RankingDataAccessor interface {
	// GetScore retrieves the score of a member in a ZSET.
	// Returns the score or 0.0 if member does not exist (along with redis.Nil error). Uses ZSCORE.
	GetScore(ctx context.Context, key string, member string) (float64, error)

	// GetTopRanked retrieves members and scores from a ZSET by rank in descending order.
	// Uses ZREVRANGE WITHSCORES.
	GetTopRanked(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)

	// Add other necessary methods like:
	// ZRank, ZRevRank, ZCard, ZRem, etc. if needed.
}

// rankingRedisDataAccessor is an implementation of RankingDataAccessor using Redis.
type rankingRedisDataAccessor struct {
	rdb *redis.Client // Redis client
}

// NewRankingDataAccessor creates a new RankingDataAccessor instance.
func NewRankingDataAccessor(rdb *redis.Client) RankingDataAccessor {
	return &rankingRedisDataAccessor{rdb: rdb}
}

// Implement GetScore using ZSCORE
func (r *rankingRedisDataAccessor) GetScore(ctx context.Context, key string, member string) (float64, error) {
	score, err := r.rdb.ZScore(ctx, key, member).Result()
	if err != nil {
		if err == redis.Nil {
			return 0.0, redis.Nil
		}
		fmt.Printf("Repo Error: ZSCORE key=%s, member=%s - %v\n", key, member, err)
		return 0.0, err
	}

	return score, nil
}

// Implement GetTopRanked using ZREVRANGE WITHSCORES
func (r *rankingRedisDataAccessor) GetTopRanked(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {

	videos, err := r.rdb.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return videos, nil
}
