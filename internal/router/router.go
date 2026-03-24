package router

import (
	"database/sql"
	"strings"

	"alioth-hrc/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func New(appEnv string, sqlDB *sql.DB, rdb *redis.Client) *gin.Engine {
	setGinMode(appEnv)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := handler.NewHealthHandler(sqlDB, rdb)
	r.GET("/healthz", h.Healthz)
	r.GET("/readyz", h.Readyz)
	r.GET("/demo/ping", h.DemoPing)

	return r
}

func setGinMode(appEnv string) {
	switch strings.ToLower(strings.TrimSpace(appEnv)) {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
