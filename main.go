package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Kaisen-dl/em-devops-task3/internal/config"
	"github.com/Kaisen-dl/em-devops-task3/internal/db"
	"github.com/Kaisen-dl/em-devops-task3/internal/handlers"
	"github.com/Kaisen-dl/em-devops-task3/internal/logger"
	"github.com/Kaisen-dl/em-devops-task3/internal/metrics"
)

func main() {
	log := logger.New()

	cfg := config.Load()
	ctx := context.Background()

	pg, err := db.NewPostgres(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Error("postgres connect failed", slog.String("err", err.Error()))
		return
	}
	defer pg.Close()

	rdb, err := db.NewRedis(ctx, cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		log.Error("redis connect failed", slog.String("err", err.Error()))
		return
	}
	defer rdb.Close()

	h := handlers.New(pg, rdb)

	mux := http.NewServeMux()
	mux.HandleFunc("/cache/set", h.CacheSet)
	mux.HandleFunc("/cache/get", h.CacheGet)
	mux.HandleFunc("/users", h.Users)
	mux.HandleFunc("/healthz", h.Healthz)
	mux.HandleFunc("/readyz", h.Readyz)
	mux.Handle("/metrics", metrics.Handler())

	// оборачиваем всё в middleware: логирование + метрики
	handler := metrics.Instrument(handlers.Logging(log)(mux))

	addr := ":" + cfg.HTTPPort
	log.Info("server starting", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Error("server stopped", slog.String("err", err.Error()))
	}
}