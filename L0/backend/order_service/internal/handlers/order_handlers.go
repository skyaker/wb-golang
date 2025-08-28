package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
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
