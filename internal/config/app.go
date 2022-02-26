package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type App struct {
	UserDb
	EventDb
	EventQueue
	Cache
	GrpcServer
}

type UserDb struct {
	Dsn    string `envconfig:"USER_DB_DSN"`
	Driver string `envconfig:"USER_DB_DRIVER"`
}

type EventDb struct {
	Dsn    string `envconfig:"EVENT_DB_DSN"`
	Driver string `envconfig:"EVENT_DB_DRIVER"`
}

type EventQueue struct {
	Broker string
	Topic  string
}

type GrpcServer struct {
	Host string
	Port string
}

type Cache struct {
	HostPort             string
	CacheExpireInSeconds time.Duration
}

func NewApp() (App, error) {
	var (
		appConfig App
		err       error
	)

	err = envconfig.Process("", &appConfig)
	if err != nil {
		return appConfig, fmt.Errorf("failed to get env variables: %w", err)
	}

	return appConfig, nil
}
