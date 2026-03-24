package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Page struct {
	Page     int
	PageSize int
	Offset   int
}

func ParsePageQuery(c *gin.Context) Page {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	ps, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if p < 1 {
		p = 1
	}
	if ps < 1 {
		ps = 20
	}
	if ps > 100 {
		ps = 100
	}
	return Page{Page: p, PageSize: ps, Offset: (p - 1) * ps}
}
