package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"alioth-hrc/internal/cache"
	"alioth-hrc/internal/config"
	"alioth-hrc/internal/db"
	"alioth-hrc/internal/model"
	"alioth-hrc/internal/router"

	"github.com/redis/go-redis/v9"
)

type App struct {
	srv   *http.Server
	sqlDB *sql.DB
	rdb   *redis.Client
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	gdb, sqlDB, err := db.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	sqlDB.SetConnMaxLifetime(cfg.MySQLConnMaxLifetime())

	rdb := cache.NewRedis(cfg.Redis)
	if err := cache.Ping(ctx, rdb); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	if err := gdb.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("mysql migrate: %w", err)
	}

	engine := router.New(cfg.AppEnv, sqlDB, rdb)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("app initialized", "name", cfg.AppName, "env", cfg.AppEnv, "port", cfg.HTTPPort)

	return &App{srv: srv, sqlDB: sqlDB, rdb: rdb}, nil
}

func (a *App) Run() error {
	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var errs []error
	errs = append(errs, a.srv.Shutdown(ctx))
	if a.sqlDB != nil {
		errs = append(errs, a.sqlDB.Close())
	}
	if a.rdb != nil {
		errs = append(errs, a.rdb.Close())
	}
	return errors.Join(errs...)
}
