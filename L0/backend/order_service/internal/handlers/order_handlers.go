package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	models "order_info/internal/models"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GetOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUIDstr := chi.URLParam(r, "order_uid")

		orderUID, err := uuid.Parse(orderUIDstr)
		if err != nil {
			log.Error().Err(err).Msg("invalid order_uid")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		orderInfo, err := getOrderInfo(db, orderUID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get order info")
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		fmt.Println(orderInfo)

		deliveryInfo, err := getDeliveryInfo(db, orderUID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get delivery info")
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		fmt.Println(deliveryInfo)

		paymentInfo, err := getPaymentInfo(db, orderUID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get payment info")
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		fmt.Println(paymentInfo)

		itemsInfo, err := getItemInfo(db, orderUID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get item info")
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		agg := struct {
			Order    models.Order    `json:"order"`
			Delivery models.Delivery `json:"delivery"`
			Payment  models.Payment  `json:"payment"`
			Items    []models.Item   `json:"items"`
		}{
			Order:    orderInfo,
			Delivery: deliveryInfo,
			Payment:  paymentInfo,
			Items:    itemsInfo,
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

func getOrderInfo(db *sql.DB, orderUID uuid.UUID) (models.Order, error) {
	var orderInfo models.Order

	query := `SELECT * FROM orders WHERE order_uid = $1`
	rows, err := db.Query(query, orderUID)
	if err != nil {
		return models.Order{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&orderInfo.OrderUID,
			&orderInfo.TrackNumber,
			&orderInfo.Entry,
			&orderInfo.Locale,
			&orderInfo.InternalSig,
			&orderInfo.CustomerID,
			&orderInfo.DeliveryServ,
			&orderInfo.ShardKey,
			&orderInfo.SMID,
			&orderInfo.DateCreated,
			&orderInfo.OOFShard,
		)
		if err != nil {
			return models.Order{}, err
		}
	}

	if rows.Err() != nil {
		return models.Order{}, rows.Err()
	}

	return orderInfo, nil
}

func getDeliveryInfo(db *sql.DB, orderUID uuid.UUID) (models.Delivery, error) {
	var deliveryInfo models.Delivery

	query := `SELECT name, phone, zip, city, address, region, email 
						FROM deliveries 
						WHERE order_uid = $1`

	rows, err := db.Query(query, orderUID)
	if err != nil {
		return models.Delivery{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&deliveryInfo.Name,
			&deliveryInfo.Phone,
			&deliveryInfo.Zip,
			&deliveryInfo.City,
			&deliveryInfo.Address,
			&deliveryInfo.Region,
			&deliveryInfo.Email,
		)
		if err != nil {
			return models.Delivery{}, err
		}
	}

	if rows.Err() != nil {
		return models.Delivery{}, rows.Err()
	}

	return deliveryInfo, nil
}

func getPaymentInfo(db *sql.DB, orderUID uuid.UUID) (models.Payment, error) {
	var paymentInfo models.Payment

	query := `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee 
						FROM payments 
						WHERE order_uid = $1`

	rows, err := db.Query(query, orderUID)
	if err != nil {
		return models.Payment{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var reqID sql.NullString

		err := rows.Scan(
			&paymentInfo.Transaction,
			&reqID,
			&paymentInfo.Currency,
			&paymentInfo.Provider,
			&paymentInfo.Amount,
			&paymentInfo.PaymentDT,
			&paymentInfo.Bank,
			&paymentInfo.DeliveryCost,
			&paymentInfo.GoodsTotal,
			&paymentInfo.CustomFee,
		)
		if err != nil {
			return models.Payment{}, err
		}

		if reqID.Valid {
			paymentInfo.RequestID = reqID.String
		}
	}

	if rows.Err() != nil {
		return models.Payment{}, rows.Err()
	}

	return paymentInfo, nil
}

func getItemInfo(db *sql.DB, orderUID uuid.UUID) ([]models.Item, error) {
	var itemsInfo []models.Item

	query := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status 
						FROM items 
						WHERE order_uid = $1`

	rows, err := db.Query(query, orderUID)
	if err != nil {
		return []models.Item{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var itemInfo models.Item
		err := rows.Scan(
			&itemInfo.ChrtID,
			&itemInfo.TrackNumber,
			&itemInfo.Price,
			&itemInfo.Rid,
			&itemInfo.Name,
			&itemInfo.Sale,
			&itemInfo.Size,
			&itemInfo.TotalPrice,
			&itemInfo.NmID,
			&itemInfo.Brand,
			&itemInfo.Status,
		)
		if err != nil {
			return []models.Item{}, err
		}
		itemsInfo = append(itemsInfo, itemInfo)
	}

	if rows.Err() != nil {
		return []models.Item{}, rows.Err()
	}

	return itemsInfo, nil
}
