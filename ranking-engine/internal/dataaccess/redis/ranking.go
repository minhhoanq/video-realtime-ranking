package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RankingDataAccessor interface {
	// IncrementScore increases the score of a member in a ZSET by a specific increment.
	// Returns the new score. Uses ZINCRBY.
	// Tăng điểm của một member trong ZSET bằng một giá trị nhất định. Trả về điểm mới. Dùng ZINCRBY.
	IncrementScore(ctx context.Context, key string, member string, increment float64) (float64, error)

	// GetScore retrieves the score of a member in a ZSET.
	// Returns the score or 0.0 if member does not exist (along with redis.Nil error). Uses ZSCORE.
	// Lấy điểm của một member trong ZSET. Trả về điểm số hoặc 0.0 nếu member không tồn tại (cùng với lỗi redis.Nil). Dùng ZSCORE.
	GetScore(ctx context.Context, key string, member string) (float64, error)

	// GetTopRanked retrieves members and scores from a ZSET by rank in descending order.
	// Uses ZREVRANGE WITHSCORES.
	// Lấy các member và điểm số từ ZSET theo thứ hạng giảm dần. Dùng ZREVRANGE WITHSCORES.
	GetTopRanked(ctx context.Context, key string, offset, limit int) ([]redis.Z, error)

	// Add other necessary methods like:
	// ZRank, ZRevRank, ZCard, ZRem, etc. if needed.
	// Thêm các phương thức khác nếu cần như ZRank, ZRevRank, ZCard, ZRem, v.v.
}

// rankingRedisDataAccessor is an implementation of RankingDataAccessor using Redis.
// rankingRedisDataAccessor là một triển khai của RankingDataAccessor sử dụng Redis.
type rankingRedisDataAccessor struct {
	rdb *redis.Client // Redis client
}

// NewRankingDataAccessor creates a new RankingDataAccessor instance.
// NewRankingDataAccessor tạo một instance RankingDataAccessor mới.
func NewRankingDataAccessor(rdb *redis.Client) RankingDataAccessor {
	return &rankingRedisDataAccessor{rdb: rdb}
}

// Implement IncrementScore using ZINCRBY
func (r *rankingRedisDataAccessor) IncrementScore(ctx context.Context, key string, member string, increment float64) (float64, error) {
	newScore, err := r.rdb.ZIncrBy(ctx, key, increment, member).Result()
	if err != nil {
		fmt.Printf("Repo Error: ZINCRBY key=%s, member=%s, inc=%f - %v\n", key, member, increment, err)
		return 0, err
	}
	return newScore, nil
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

	// Parse string score sang float
	if err != nil {
		return 0.0, err
	}

	return score, nil
}

// Implement GetTopRanked using ZREVRANGE WITHSCORES
func (r *rankingRedisDataAccessor) GetTopRanked(ctx context.Context, key string, offset, limit int) ([]redis.Z, error) {
	start := int64(offset)
	stop := int64(offset + limit - 1)

	// ZRevRangeWithScores lấy theo thứ tự giảm dần điểm số
	videos, err := r.rdb.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		fmt.Printf("Repo Error: ZREVRANGE key=%s, offset=%d, limit=%d - %v\n", key, offset, limit, err)
		return nil, err
	}
	return videos, nil
}
