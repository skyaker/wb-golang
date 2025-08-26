package main

import (
	"net/http"

	handlers "order_info/internal/handlers"
	order_kafka "order_info/internal/kafka"
	rep "order_info/internal/repository"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func main() {
	db := rep.GetDbConnection()
	defer db.Close()

	go order_kafka.RunKafkaListener(db)

	r := chi.NewRouter()

	r.Get("/orders/{order_uid}", handlers.GetOrder(db))

	log.Info().Msg("Order server is running")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Server start failed")
	}
}
