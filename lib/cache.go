package lib

import (
	"fmt"
	"time"

	"fampay-assignment/config"
	"fampay-assignment/connections"
	"fampay-assignment/logger"
	"fampay-assignment/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func GetCache(key string) (string, error) {
	ctx := connections.GetContext()

	val, err := connections.RedisClient.Get(ctx, key).Result()
	return val, err
}

func SetCache(key string, val []byte, ttl time.Duration) error {
	ctx := connections.GetContext()
	return connections.RedisClient.Set(ctx, key, val, ttl).Err()
}

func createCacheKey( op string, params any) (key string, err error) {
	encodedParams, err := utils.EncodeToGob(params)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"params": params,
			"err":    err,
		}).Error("failed to encode cache params")
		return key, err
	}

	return fmt.Sprintf("videos:%s:%s", op, encodedParams), nil
}

func execAndCacheQueryResult[Params any, Result any](
	key string,
	query func(*pgxpool.Pool, *Params) Result,
	db *pgxpool.Pool,
	params *Params,
) (result Result) {
	result = query(db, params)

	resultEncoded, err := utils.EncodeToGob(result)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"result": result,
			"err":    err,
		}).Error("failed to encode query result")
		return result
	}

	err = SetCache(key, resultEncoded, config.CACHE_TTL)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"key": key,
			"err": err,
		}).Error("failed to set cache")
	}

	return result
}

func cacheQuery[Params any, Result any](
	name string,
	query func(*pgxpool.Pool, *Params) Result,
	db *pgxpool.Pool,
	params *Params,
) (result Result) {
	key, err := createCacheKey(name, params)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"query":  name,
			"params": params,
			"err":    err,
		}).Error("failed to create cache key")
		return query(db, params)
	}

	cacheResult, err := GetCache(key)
	if err != nil {
		if err == redis.Nil {
			logger.Log.WithFields(logger.Fields{
				"query": name,
				"key":   key,
			}).Debug("cache miss")
		} else {
			logger.Log.WithFields(logger.Fields{
				"query": name,
				"key":   key,
				"err":   err,
			}).Error("failed to get from cache")
		}
		return execAndCacheQueryResult(key, query, db, params)
	}

	err = utils.DecodeFromGob([]byte(cacheResult), &result)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"key": key,
			"err": err,
		}).Error("failed to decode cached result")
		return execAndCacheQueryResult(key, query, db, params)
	}

	logger.Log.WithFields(logger.Fields{
		"query": name,
		"key":   key,
	}).Info("cache hit")
	return result
}

func InvalidateCache(org string) (keysDeleted int64, err error) {
	ctx := connections.GetContext()

	const batchSize = 10000
	var pattern string
	if org == "" {
		pattern = "dashboard:*"
	} else {
		pattern = fmt.Sprintf("dashboard:%s:*", org)
	}

	var cursor uint64
	for {
		keys, cursor, err := connections.RedisClient.Scan(
			ctx,
			cursor,
			pattern,
			batchSize,
		).Result()
		if err != nil {
			logger.Log.WithFields(logger.Fields{
				"org": org,
				"err": err,
			}).Error("error scanning keys")
			return keysDeleted, err
		}

		if len(keys) == 0 {
			break
		}

		deleted, err := connections.RedisClient.Unlink(
			ctx,
			keys...,
		).Result()
		if err != nil {
			logger.Log.WithFields(logger.Fields{
				"org":   org,
				"count": deleted,
				"err":   err,
			}).Error("error unlinking keys")
		}

		logger.Log.WithFields(logger.Fields{
			"org":   org,
			"count": deleted,
		}).Info("unlinked keys")
		keysDeleted += deleted

		if cursor == 0 {
			break
		}
	}

	logger.Log.WithFields(logger.Fields{
		"org":   org,
		"count": keysDeleted,
	}).Info("invalidated cache")
	return keysDeleted, nil
}
