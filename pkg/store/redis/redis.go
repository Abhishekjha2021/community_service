package redis

import (
	"context"
	"fmt"

	"github.com/Abhishekjha321/community_service/storage/cache"
	"github.com/Abhishekjha321/community_service/pkg/config"
)

var redisCluster cache.CacheBase

func InitializeRedis() (cache.CacheBase, error) {
	var err error
	ctx := context.Background()
	redisCluster, err = cache.NewRedisClient(ctx, config.Config.RedisConfig.HostNames, config.Config.RedisConfig.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}
	return redisCluster, nil
}

func Client() cache.CacheBase {
	return redisCluster
}
