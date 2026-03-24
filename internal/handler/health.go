package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"alioth-hrc/internal/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	sqlDB *sql.DB
	rdb   *redis.Client
}

func NewHealthHandler(sqlDB *sql.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{sqlDB: sqlDB, rdb: rdb}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.sqlDB.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "mysql": err.Error()})
		return
	}

	if err := cache.Ping(ctx, h.rdb); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "redis": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func (h *HealthHandler) DemoPing(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var mysqlOne int
	if err := h.sqlDB.QueryRowContext(ctx, "SELECT 1").Scan(&mysqlOne); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "fail", "mysql": err.Error()})
		return
	}

	const key = "demo:ping"
	const val = "pong"
	if err := h.rdb.Set(ctx, key, val, 10*time.Second).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "fail", "redis": err.Error()})
		return
	}
	got, err := h.rdb.Get(ctx, key).Result()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "fail", "redis": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"mysql":  "ok",
		"redis":  "ok",
		"result": gin.H{
			"mysql_select": mysqlOne,
			"redis_value":  got,
		},
	})
}
