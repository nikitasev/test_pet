package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type EventLog struct {
	Db *sql.DB
}

func NewEventLogStorage(db *sql.DB) *EventLog {
	return &EventLog{Db: db}
}

func (l *EventLog) Log(batch []LogMessage) error {
	tx, err := l.Db.Begin()
	if err != nil {
		return err
	}

	valueStrings := make([]string, 0, len(batch))
	valueArgs := make([]interface{}, 0, len(batch)*2)
	for _, msg := range batch {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, msg.date)
		valueArgs = append(valueArgs, msg.userId)
	}

	stmt, err := tx.Prepare(fmt.Sprintf(`INSERT INTO log (date_time, user_id) VALUES %s`, strings.Join(valueStrings, ",")))
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(valueArgs...); err != nil {
		return err
	}

	return tx.Commit()
}

type LogMessage struct {
	userId int64
	date   time.Time
}

func NewLogMessage(userId int64, date time.Time) LogMessage {
	return LogMessage{userId: userId, date: date}
}
