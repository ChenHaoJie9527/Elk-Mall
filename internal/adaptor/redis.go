package adaptor

import (
	"context"
	"fmt"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
	"github.com/redis/go-redis/v9"
)

// 打开 Redis 连接
func OpenRedis(cfg config.RedisConfig) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := PingRedis(context.Background(), r); err != nil {
		return nil, fmt.Errorf("ping Redis 失败: %w", err)
	}

	return r, nil
}

// ping Redis 连接
func PingRedis(ctx context.Context, r *redis.Client) error {
	if err := r.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping Redis 失败: %w", err)
	}
	return nil
}
