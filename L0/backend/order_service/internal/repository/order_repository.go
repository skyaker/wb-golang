package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	models "order_info/internal/models"
)

func WriteOrder(db *sql.DB, order *models.Order) error {
	log.Info().Msg("Write order")
	return nil
}
