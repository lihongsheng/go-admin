// Package initialize 初始化
package initialize

import (
	"context"
	"fmt"

	"go-admin/server/global"

	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化 Redis 客户端
func InitRedis() {
	cfg := global.Cfg.Redis
	if !cfg.Enable {
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("redis ping failed: %v", err))
	}
	global.Redis = rdb
}
