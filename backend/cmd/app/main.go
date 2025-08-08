package main

import (
	"context"
	"github.com/ilam072/wbtech-l0/backend/internal/broker/kafka/consumer"
	"github.com/ilam072/wbtech-l0/backend/internal/broker/kafka/handler"
	"github.com/ilam072/wbtech-l0/backend/internal/cache"
	"github.com/ilam072/wbtech-l0/backend/internal/config"
	"github.com/ilam072/wbtech-l0/backend/internal/converter"
	"github.com/ilam072/wbtech-l0/backend/internal/repo/postgres"
	"github.com/ilam072/wbtech-l0/backend/internal/rest"
	"github.com/ilam072/wbtech-l0/backend/internal/service"
	"github.com/ilam072/wbtech-l0/backend/internal/validator"
	"github.com/ilam072/wbtech-l0/backend/pkg/db"
	"github.com/ilam072/wbtech-l0/backend/pkg/logger/handlers/slogpretty"
	"github.com/ilam072/wbtech-l0/backend/pkg/logger/sl"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title Order Service
// @description REST API service using Kafka, PostgreSQL and in-memory cache

// @host localhost:8082
// @BasePath /
func main() {
	cfg := config.New()

	pool, err := db.OpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	l := initLogger()

	orderValidator := validator.New()
	orderRepo := postgres.NewOrderRepo(pool)
	converterr := converter.New()
	cache := cache.New(orderRepo, converterr)
	orderService := service.NewOrderService(orderRepo, cache, converterr)

	kafkaConsumer := consumer.New(
		cfg.KafkaConfig.Topic,
		cfg.KafkaConfig.GroupID,
		cfg.KafkaConfig.Brokers...,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	orderConsumerHandler := handler.NewOrderConsumerHandler(
		l,
		kafkaConsumer,
		orderService,
		orderValidator,
	)

	go func() {
		if err := orderConsumerHandler.Start(ctx); err != nil {
			l.Error("failed to start consumer", sl.Err(err))
		}
	}()

	err = cache.Preload(context.Background(), cfg.CacheConfig.PreloadLimit)
	if err != nil {
		log.Fatalln("error preloading cache", sl.Err(err))
	}

	h := rest.NewHandler(l, orderService)
	go func() {
		if err := h.Listen(cfg.ServerConfig.Address()); err != nil {
			l.Error("failed to start server", sl.Err(err))
			cancel()
		}
	}()

	<-sigs
	l.Info("shutting down...")
	cancel()

	if err := h.Shutdown(); err != nil {
		l.Error("failed to shutdown server", sl.Err(err))
	}

	if err := kafkaConsumer.Close(); err != nil {
		l.Error("failed to close kafka consumer", sl.Err(err))
	}

	l.Info("application stopped")
}

func initLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	h := opts.NewPrettyHandler(os.Stdout)

	return slog.New(h)
}
