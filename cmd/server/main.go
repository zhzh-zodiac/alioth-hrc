// Package main 启动 HTTP 服务。Swagger 元信息见下方注释（供 swag 生成文档）。
//
//	@title			alioth-hrc API
//	@version		1.0
//	@description	人情往来记账后端：认证、联系人、账本、礼金分类与流水、统计等。
//	@BasePath		/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				访问令牌，格式：Bearer {access_token}
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "alioth-hrc/docs"

	"alioth-hrc/internal/app"
)

//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g main.go -o ../../docs --parseDependency --parseInternal

func main() {
	ctx := context.Background()
	a, err := app.New(ctx)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := a.Run(); err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
