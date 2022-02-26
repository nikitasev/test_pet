package service

import (
	"github.com/segmentio/kafka-go"
	"test_pet/internal/config"
)

func NewKafkaConsumer(cfg config.EventQueue) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{cfg.Broker},
		Topic:     cfg.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
}
