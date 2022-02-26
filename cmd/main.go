package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"
	"test_pet/internal/application"
	"test_pet/internal/config"
	"test_pet/internal/infrastructure/persistence"
	"test_pet/internal/infrastructure/service"
	"test_pet/pkg/grpc/userapi"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg, err := config.NewApp()
	if err != nil {
		panic(err)
	}

	userDb, err := service.Connect(cfg.UserDb.Dsn, cfg.UserDb.Driver)
	if err != nil {
		panic(err)
	}
	defer userDb.Close()
	userRepo := persistence.NewUserRepository(userDb)

	eventDb, err := service.Connect(cfg.EventDb.Dsn, cfg.EventDb.Driver)
	if err != nil {
		panic(err)
	}
	defer eventDb.Close()
	eventRepo := persistence.NewEventLogStorage(eventDb)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GrpcServer.Host, cfg.GrpcServer.Port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	client := redis.NewClient(&redis.Options{Addr: cfg.HostPort})
	if err := client.Ping(); err != nil {
		panic(err)
	}
	defer client.Close()
	userListCache := service.NewRedisUserListCache(client, cfg.Cache)

	kafkaProducer := service.NewKafkaProducer(cfg.EventQueue)
	defer kafkaProducer.Close()
	eventLogger := service.NewQueueEventLogger(kafkaProducer)

	kafkaConsumer := service.NewKafkaConsumer(cfg.EventQueue)
	defer kafkaConsumer.Close()

	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	userapi.RegisterUserServiceServer(grpcServer, application.NewGrpcHandler(userRepo, eventLogger, userListCache, logger))

	eventQueueHandler := application.NewEventQueueHandler(kafkaConsumer, eventRepo)

	queueStop := make(chan bool, 1)
	grpcStop := make(chan error, 1)
	osStop := make(chan os.Signal, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		eventQueueHandler.HandleQueue(queueStop)
		wg.Done()
	}()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error(err.Error())
			grpcStop <- err
		}
		wg.Done()
	}()

	func() {
		for {
			select {
			//If OS terminating process: gracefully stopping GRPC server and Event Queue handler
			case <-osStop:
				queueStop <- true
				grpcServer.GracefulStop()
				return
			case err := <-grpcStop:
				logger.Error("GRPC server error", zap.Error(err))
				queueStop <- true
				return
			default:
				continue
			}
		}
	}()

	//Waiting gracefully stop of all routines
	wg.Wait()
	logger.Info("terminating ...")

	os.Exit(0)
}
