package service

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

type QueueEventLogger struct {
	producer *kafka.Writer
}

type QueueItem struct {
	UserId int64     `json:"user_id" `
	Time   time.Time `json:"time"`
}

func NewQueueEventLogger(producer *kafka.Writer) *QueueEventLogger {
	return &QueueEventLogger{producer: producer}
}

func (l *QueueEventLogger) Log(userId int64) error {
	msg := QueueItem{UserId: userId, Time: time.Now()}
	b, _ := json.Marshal(msg)

	return l.producer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("new_user"),
			Value: b,
		},
	)
}
