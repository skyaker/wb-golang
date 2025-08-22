package repository

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	models "order_info/internal/models"
)

func WriteOrder(db *sql.DB, order *models.Order) error {
	order.OrderUID = uuid.New().String()

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
		reqID = nil
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
