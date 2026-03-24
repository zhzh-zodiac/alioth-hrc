package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func queryUintPtr(c *gin.Context, key string) *uint {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	u, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil
	}
	v := uint(u)
	return &v
}

func queryIntPtr(c *gin.Context, key string) *int {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

func queryDatePtr(c *gin.Context, key string) (*time.Time, error) {
	s := c.Query(key)
	if s == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
