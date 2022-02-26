package service

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(dsn, driver string) (*sqlx.DB, error) {
	return sqlx.Connect(driver, dsn)
}
