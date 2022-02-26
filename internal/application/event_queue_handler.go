package application

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"test_pet/internal/infrastructure/persistence"
	"test_pet/internal/infrastructure/service"
)

type EventQueueHandler struct {
	consumer        *kafka.Reader
	eventLogStorage *persistence.EventLog
	logger          *zap.Logger
}

func NewEventQueueHandler(consumer *kafka.Reader, eventLogStorage *persistence.EventLog) *EventQueueHandler {
	return &EventQueueHandler{
		consumer:        consumer,
		eventLogStorage: eventLogStorage,
	}
}

func (h *EventQueueHandler) HandleQueue(cStop <-chan bool) {
	for {
		select {
		case <-cStop:
			return
		default:
			m, err := h.consumer.ReadMessage(context.Background())
			if err != nil {
				h.logger.Error("failed reading message from queue", zap.Error(err))
				continue
			}
			var msg service.Message
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				h.logger.Error("failed encoding message", zap.Error(err))
				continue
			}
			if err := h.eventLogStorage.Log(msg.UserId, msg.Time); err != nil {
				h.logger.Error("failed logging event", zap.Error(err))
			}
		}
	}
}
