package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	cache "order_info/internal/cache"
	models "order_info/internal/models"
	rep "order_info/internal/repository"
	"time"

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

		cacheStart := time.Now()
		ok, err := cache.GetOrder(orderUID.String(), &agg)
		if err != nil {
			log.Info().Err(err).Msg("failed to get order from cache")
		}
		if ok {
			cacheEnd := time.Now()
			log.Info().
				Str("order_uid", orderUID.String()).
				Dur("time", cacheEnd.Sub(cacheStart)).
				Msg("get order from cache")
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

		readStart := time.Now()
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
		readEnd := time.Now()
		log.Info().
			Str("order_uid", orderUID.String()).
			Dur("time", readEnd.Sub(readStart)).
			Msg("read order from db")

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
