package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func main() {
	r := chi.NewRouter()

	log.Info().Msg("Order server is running")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth service").
			Msg("Server start failed")
	}
}
