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
	Dsn string `envconfig:"USER_DB_DSN"`
}

type EventDb struct {
	Dsn string `envconfig:"EVENT_DB_DSN"`
}

type EventQueue struct {
	Broker string `envconfig:"EVENT_QUEUE_BROKER"`
	Topic  string `envconfig:"EVENT_QUEUE_TOPIC"`
}

type GrpcServer struct {
	HostPort string `envconfig:"GRPC_HOST_PORT"`
}

type Cache struct {
	HostPort             string        `envconfig:"CACHE_HOST_PORT"`
	CacheExpireInSeconds time.Duration `envconfig:"CACHE_EXPIRE_IN_SECONDS"`
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
