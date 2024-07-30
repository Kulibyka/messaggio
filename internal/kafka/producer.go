package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"messagio/internal/config"

	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: writer}, nil
}

func (p *Producer) WriteMessage(key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	err := p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
