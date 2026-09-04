package adaptor

import (
	"context"
	"testing"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
)

func TestOpenRedis(t *testing.T) {
	cfg := config.RedisConfig{
		Addr:     "127.0.0.1:6380",
		Password: "",
		DB:       0,
	}

	r, err := OpenRedis(cfg)
	if err != nil {
		t.Fatalf("打开 Redis 失败: %v", err)
	}

	defer r.Close()

	if err := PingRedis(context.Background(), r); err != nil {
		t.Fatalf("ping Redis 失败: %v", err)
	}
}
