package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	models "order_info/internal/models"
	rep "order_info/internal/repository"

	validator "github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	kafka "github.com/segmentio/kafka-go"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func RunKafkaListener(db *sql.DB) {
	kafkaURL := fmt.Sprintf("%v:%v", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	topics := []string{"add-order"}
	groupID := "1"

	reader := getKafkaReader(kafkaURL, topics, groupID)

	defer reader.Close()

	log.Info().Msg("Start consuming kafka topic")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("kafka message reading error")
			continue
		}

		switch m.Topic {
		case "add-order":
			var order models.Order

			if err := json.Unmarshal(m.Value, &order); err != nil {
				log.Error().Err(err).Msg("json parse error")
				continue
			}
			log.Info().Msg("add-order message received")

			if err := validateOrder(&order); err != nil {
				log.Warn().Err(err).Msg("validation failed")
				continue
			}

			if err := rep.WriteOrder(db, &order); err != nil {
				log.Error().Err(err).Msg("failed to save order")
				continue
			}
			log.Info().Msg("Order created")

		default:
			log.Error().Msg("topic undefined")
		}
	}
}

func getKafkaReader(kafkaURL string, topics []string, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})
}

func validateOrder(order *models.Order) error {
	if err := validate.Struct(order); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			var b strings.Builder
			for _, fe := range verrs {
				b.WriteString(fmt.Sprintf("Field %s: invalid '%s'; ", fe.Field(), fe.Tag()))
			}
			return fmt.Errorf("validation errors: %s", b.String())
		}
		return fmt.Errorf("Validation process error: %v", err)
	}
	return nil
}
