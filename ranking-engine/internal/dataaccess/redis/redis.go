package redis

import (
	"context"
	"fmt"
	"video-realtime-ranking/ranking-engine/config"

	"github.com/redis/go-redis/v9"
)

const (
	GlobalScore    = "video:global_score"
	UserVideoScore = "video"
)

type Redis struct {
	cfg config.Config
}

func NewRedis(cfg config.Config) *Redis {
	return &Redis{
		cfg: cfg,
	}
}

func (r *Redis) Connect() (*redis.Client, error) {
	// connect to redis
	otps := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", r.cfg.Redis.Host, r.cfg.Redis.Port),
	}

	client := redis.NewClient(otps)
	// check connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot connection to redis, err: ", err.Error())
	}
	fmt.Println("connect to redis successfully")

	return client, nil
}
