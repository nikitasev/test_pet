package application

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"test_pet/internal/infrastructure/persistence"
	"test_pet/internal/infrastructure/service"
	"time"
)

const heapMaxSize = 1000

type EventQueueHandler struct {
	heap            []service.Message
	consumer        *kafka.Reader
	eventLogStorage *persistence.EventLog
	logger          *zap.Logger
}

func NewEventQueueHandler(consumer *kafka.Reader, eventLogStorage *persistence.EventLog, logger *zap.Logger) *EventQueueHandler {
	return &EventQueueHandler{
		consumer:        consumer,
		eventLogStorage: eventLogStorage,
		logger:          logger,
	}
}

func (h *EventQueueHandler) HandleQueue(cStop <-chan bool) {
	for {
		select {
		case <-cStop:
			return
		default:
			h.handle()
		}
	}
}

func (h *EventQueueHandler) handle() {
	d := time.Now().Add(time.Second * 5)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	m, err := h.consumer.ReadMessage(ctx)
	if err != nil {
		if !errors.Is(context.DeadlineExceeded, err) {
			h.logger.Error("failed reading message from queue", zap.Error(err))
		}
		return
	}
	var msg service.Message
	if err := json.Unmarshal(m.Value, &msg); err != nil {
		h.logger.Error("failed encoding message", zap.Error(err))
		return
	}
	if h.heap == nil {
		h.heap = make([]service.Message, 0)
	}
	h.heap = append(h.heap, msg)
	if len(h.heap) >= heapMaxSize {
		if err := h.eventLogStorage.Log(msg.UserId, msg.Time); err != nil {
			h.logger.Error("failed logging event", zap.Error(err))
		}
	}
}
