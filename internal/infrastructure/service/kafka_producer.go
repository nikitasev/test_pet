package service

import (
	"github.com/segmentio/kafka-go"
	"test_pet/internal/config"
)

func NewKafkaProducer(cfg config.EventQueue) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(cfg.Broker),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}
