package connections

import (
	"context"

	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"

	"fampay-assignment/config"
	"fampay-assignment/logger"
	"fampay-assignment/utils"
)

var (
	RedisClient *redis.Client
)

func GetRedisContext() context.Context {
	return utils.GetContextWithTimeout(ctx, config.REDIS_TIMEOUT)
}

func connectRedis(
	uri string,
) *redis.Client {
	redisOptions, err := redis.ParseURL(uri)
	if err != nil {
		logger.Log.WithField(
			"err", err,
		).Fatal("failed to connect to redis")
	}

	client := redis.NewClient(redisOptions)
	client.AddHook(nrredis.NewHook(redisOptions))

	logger.Log.Info("connected to redis")
	return client
}

func init() {
	RedisClient = connectRedis(config.RedisUri)
}
