package repository

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	models "order_info/internal/models"
)

type HttpError struct {
	Code  int    `json:"code"`
	Error error  `json:"error"`
	Msg   string `json:"msg"`
}

// Inserts all info into the database
func WriteOrder(db *sql.DB, order *models.Order) error {
	_, err := uuid.Parse(order.OrderUID)
	if err != nil {
		log.Info().Msg("Order UID parse error")
		order.OrderUID = uuid.New().String()
	}

	if err := insertOrderInfo(db, order); err != nil {
		log.Error().
			Err(err).
			Msg("Order info write failed")
		return err
	}

	if err := insertDeliveryInfo(db, order); err != nil {
		log.Error().
			Err(err).
			Msg("Delivery info write failed")
		return err
	}

	if err := insertPaymentInfo(db, order); err != nil {
		log.Error().
			Err(err).
			Msg("Payment info write failed")
		return err
	}

	if err := insertItemInfo(db, order); err != nil {
		log.Error().
			Err(err).
			Msg("Item info write failed")
		return err
	}

	log.Info().Msg("Order write success")
	return nil
}

func insertOrderInfo(db *sql.DB, order *models.Order) error {
	orderQuery := `
    INSERT INTO orders (
    	order_uid,
      track_number,
      entry,
      locale,
      internal_signature,
      customer_id,
      delivery_service,
      shardkey,
      sm_id,
      date_created,
      oof_shard
    )
    VALUES (
      $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
    )`

	_, err := db.Exec(
		orderQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSig,
		order.CustomerID,
		order.DeliveryServ,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OOFShard,
	)
	if err != nil {
		return err
	}
	return nil
}

func insertDeliveryInfo(db *sql.DB, order *models.Order) error {
	deliveryQuery := `
    INSERT INTO deliveries (
      order_uid,
      name,
      phone,
      zip,
      city,
      address,
      region,
      email
    )
    VALUES (
      $1, $2, $3, $4, $5, $6, $7, $8
    )`

	_, err := db.Exec(
		deliveryQuery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return err
	}
	return nil
}

func insertPaymentInfo(db *sql.DB, order *models.Order) error {
	var reqID any
	if order.Payment.RequestID == "" {
		reqID = nil // save empty string as null for error avoidance
	} else {
		reqID = order.Payment.RequestID
	}

	paymentQuery := `
    INSERT INTO payments (
      order_uid,
      transaction,
      request_id,
      currency,
      provider,
      amount,
      payment_dt,
      bank,
      delivery_cost,
      goods_total,
      custom_fee
    )
    VALUES (
      $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
    )`

	_, err := db.Exec(
		paymentQuery,
		order.OrderUID,
		order.Payment.Transaction,
		reqID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}
	return nil
}

func insertItemInfo(db *sql.DB, order *models.Order) error {
	for _, item := range order.Items {
		itemQuery := `
      INSERT INTO items (
        order_uid,
        chrt_id,
        track_number,
        price,
        rid,
        name,
        sale,
        size,
        total_price,
        nm_id,
        brand,
        status
      )
      VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
      )`

		_, err := db.Exec(
			itemQuery,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadOrder(db *sql.DB, orderUID uuid.UUID) (models.AggregatedOrder, HttpError) {
	var httpErr HttpError
	orderInfo, err := getOrderInfo(db, orderUID)
	if err != nil {
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "failed to get order info"
		return models.AggregatedOrder{}, httpErr
	}

	deliveryInfo, err := getDeliveryInfo(db, orderUID)
	if err != nil {
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "failed to get delivery info"
		return models.AggregatedOrder{}, httpErr
	}

	paymentInfo, err := getPaymentInfo(db, orderUID)
	if err != nil {
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "failed to get payment info"
		return models.AggregatedOrder{}, httpErr
	}

	itemsInfo, err := getItemInfo(db, orderUID)
	if err != nil {
		httpErr.Code = http.StatusInternalServerError
		httpErr.Error = err
		httpErr.Msg = "failed to get item info"
		return models.AggregatedOrder{}, httpErr
	}

	// Aggregate order info into a single struct
	agg := models.AggregatedOrder{
		Order:    orderInfo,
		Delivery: deliveryInfo,
		Payment:  paymentInfo,
		Items:    itemsInfo,
	}
	return agg, HttpError{}
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
		// Temp null string var to avoid error of missing unrequired req id
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
