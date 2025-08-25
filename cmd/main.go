//go:build !windows
// +build !windows

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paudarco/doc-storage/internal/cache"
	"github.com/paudarco/doc-storage/internal/config"
	"github.com/paudarco/doc-storage/internal/handler"
	"github.com/paudarco/doc-storage/internal/repository"
	"github.com/paudarco/doc-storage/internal/service"
	"github.com/paudarco/doc-storage/pkg/logger"
	"github.com/paudarco/doc-storage/pkg/postgres"
	"github.com/paudarco/doc-storage/pkg/redis"
	"github.com/paudarco/doc-storage/pkg/server"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.InitLogger(cfg.Env)

	// Create connection pool to db
	pool, err := postgres.NewPostgresPool(cfg.DB)
	if err != nil {
		log.Fatalf("error creating pool: %s", err.Error())
	}
	defer pool.Close()

	redis, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("error connecting redis client: %s", err.Error())
	}

	repos := repository.NewRepository(pool)
	cache := cache.NewCache(redis, cfg)
	services := service.NewService(repos, cache, cfg, log)
	handler := handler.NewHandler(services, cfg, log)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.Server, handler.InitRoutes()); err != nil {
			log.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Println("docs storage started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Stopping docs storage...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("error while shutting down: %s\n", err.Error())
	}

	pool.Close()

	log.Println("server was stopped")
}
