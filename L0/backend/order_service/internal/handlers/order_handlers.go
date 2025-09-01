package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	cache "order_info/internal/cache"
	models "order_info/internal/models"
	rep "order_info/internal/repository"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Handler function for GET /orders/{order_uid}
// Returns order info in order_service/model.json format
func GetOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var httpErr rep.HttpError
		var agg models.AggregatedOrder

		orderUIDstr := chi.URLParam(r, "order_uid")

		orderUID, err := uuid.Parse(orderUIDstr)
		if err != nil {
			log.Error().Err(err).Msg("invalid order_uid")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ok, err := cache.GetOrder(orderUID.String(), &agg)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get order from cache")
		}
		if ok {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(agg); err != nil {
				log.Error().Err(err).Msg("failed to encode response")
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
				return
			}
			return
		}

		agg, httpErr = rep.ReadOrder(db, orderUID)
		if httpErr.Error != nil {
			log.Error().Err(httpErr.Error).Msg(httpErr.Msg)
			http.Error(
				w,
				httpErr.Msg,
				httpErr.Code,
			)
			return
		}

		cacheOrder := agg.Order
		cacheOrder.Delivery = &agg.Delivery
		cacheOrder.Payment = &agg.Payment
		cacheOrder.Items = agg.Items

		cacheData, err := json.Marshal(cacheOrder)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal order info")
		}

		err = cache.SetOrder(context.Background(), orderUID.String(), cacheData)
		if err != nil {
			key := "order:" + orderUID.String()
			log.Error().Err(err).Str("key", key).Msg("failed to set order in cache")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(agg); err != nil {
			log.Error().Err(err).Msg("failed to encode response")
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
	}
}
