package service

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mailru/go-clickhouse"
)

func PostgresConnect(dsn string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", dsn)
}

func ClickHouseConnect(dsn string) (*sql.DB, error) {
	connect, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	if err := connect.Ping(); err != nil {
		return nil, err
	}
	return sql.Open("clickhouse", dsn)
}
