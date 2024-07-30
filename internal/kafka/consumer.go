package kafka

import (
	"context"
	"database/sql"
	"log"
	"messagio/internal/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	db     *sql.DB
}

func NewConsumer(cfg config.KafkaConfig, db *sql.DB) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   cfg.Brokers,
		Topic:     cfg.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	return &Consumer{reader: reader, db: db}
}

func (c *Consumer) ConsumeMessages() {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message from Kafka: %v", err)
			break
		}

		log.Printf("Received message: %s", string(m.Value))

		_, err = c.db.Exec("UPDATE messages SET processed = TRUE, processed_at = $1 WHERE content = $2",
			time.Now(), string(m.Value))
		if err != nil {
			log.Printf("Failed to update message status: %v", err)
			break
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
