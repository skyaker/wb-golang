package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	rdb *redis.Client
	ctx = context.Background()
	ttl = time.Duration(24 * time.Hour)
)

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("Redis connection failed")
	}
	log.Info().Msg("Redis connection successful")
}

func SetOrder(ctx context.Context, orderUID string, data []byte) error {
	key := "order:" + orderUID
	if err := rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("failed to set order in cache")
		return err
	}
	return nil
}

func GetOrder(orderUID string, target any) (bool, error) {
	val, err := rdb.Get(ctx, "order:"+orderUID).Result()
	if err == redis.Nil {
		return false, nil // No order info in cache
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), target); err != nil {
		return false, err
	}
	return true, nil
}
