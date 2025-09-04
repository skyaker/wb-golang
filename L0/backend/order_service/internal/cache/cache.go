package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"order_info/internal/models"
	rep "order_info/internal/repository"
	"os"
	"time"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
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

func CacheWarmUp(db *sql.DB) {
	log.Info().Msg("Cache warm up started")

	var orderUIDs []uuid.UUID

	orderUIDs, err := getFreshOrdersId(db)
	if err != nil {
		log.Error().Err(err).Msg("Cache warm up: failed to get fresh orders")
		return
	}

	for _, orderUID := range orderUIDs {
		var tempOrder models.Order

		orderInfo, httpErr := rep.ReadOrder(db, orderUID)
		if httpErr.Error != nil {
			log.Error().Err(httpErr.Error).Msg(httpErr.Msg)
			continue
		}

		tempOrder = orderInfo.Order
		tempOrder.Delivery = &orderInfo.Delivery
		tempOrder.Payment = &orderInfo.Payment
		tempOrder.Items = orderInfo.Items

		var orderJSON []byte
		orderJSON, err = json.Marshal(tempOrder)
		if err != nil {
			log.Error().Err(err).Msg("Cache warm up: failed to marshal order info")
			continue
		}

		err = SetOrder(context.Background(), orderUID.String(), orderJSON)
		if err != nil {
			log.Error().Err(err).Msg("Cache warm up: failed to set order in cache")
			continue
		}
	}

	log.Info().Msg("Cache warm up: done")
}

func getFreshOrdersId(db *sql.DB) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	query := `SELECT order_uid FROM orders WHERE date_created > $1`
	rows, err := db.Query(query, time.Now().AddDate(0, 0, -1))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
