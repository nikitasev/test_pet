package persistence

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type EventLog struct {
	Db *sqlx.DB
}

func NewEventLogStorage(db *sqlx.DB) *EventLog {
	return &EventLog{Db: db}
}

func (l *EventLog) Log(userId int64, date time.Time) error {
	return nil
}
