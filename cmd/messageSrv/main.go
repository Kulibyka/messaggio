package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"messagio/internal/config"
	api "messagio/internal/http"
	"messagio/internal/kafka"
	"messagio/internal/storage/postgresql"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("app started")

	db, err := postgresql.New(cfg.PostgreSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	consumer := kafka.NewConsumer(cfg.Kafka, db.GetDB())
	defer consumer.Close()

	go consumer.ConsumeMessages()

	api.RegisterRoutes(router, db, producer)

	addr := ":8080"
	log.Info("Starting server", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
