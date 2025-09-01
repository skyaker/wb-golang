package cache

import (
	"context"
	"encoding/json"
	"order_info/internal/models"
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
		return err
	}
	return nil
}

func GetOrder(orderUID string, target *models.AggregatedOrder) (bool, error) {
	val, err := rdb.Get(ctx, "order:"+orderUID).Result()
	var order models.Order
	if err == redis.Nil {
		return false, nil // No order info in cache
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return false, err
	}
	*target = models.AggregatedOrder{
		Order:    order,
		Delivery: *order.Delivery,
		Payment:  *order.Payment,
		Items:    order.Items,
	}

	return true, nil
}
