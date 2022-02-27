package persistence

import (
	"database/sql"
	"time"
)

type EventLog struct {
	Db *sql.DB
}

func NewEventLogStorage(db *sql.DB) *EventLog {
	return &EventLog{Db: db}
}

func (l *EventLog) Log(userId int64, date time.Time) error {
	tx, err := l.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO log (date_time, user_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(date, userId); err != nil {
		return err
	}

	return tx.Commit()
}
