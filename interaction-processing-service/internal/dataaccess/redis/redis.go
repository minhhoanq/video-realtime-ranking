package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/minhhoanq/video-realtime-ranking/interaction-processing-service/config"

	"github.com/redis/go-redis/v9"
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

func (r *Redis) AcquireLock(client *redis.Client, lockKey string, lockValue string, timeout time.Duration) bool {
	ctx := context.Background()

	// Try to acquire the lock with SETNX command (SET if Not Exists)
	lockAcquire, err := client.SetNX(ctx, lockKey, lockValue, timeout).Result()
	if err != nil {
		fmt.Println("error acquiring lock: ", err.Error())
		return false
	}

	return lockAcquire
}

func (r *Redis) ReleaseLock(client *redis.Client, lockKey string) error {
	ctx := context.Background()
	_, err := client.Del(ctx, lockKey).Result()
	if err != nil {
		return err
	}

	return nil
}
