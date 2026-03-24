package handler

import (
	"net/http"

	"alioth-hrc/internal/middleware"
	"alioth-hrc/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	svc *service.StatsService
}

func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

// ContactSummaries 按联系人汇总收/送差额等。
//
//	@Summary		按联系人统计
//	@Tags			stats
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	handler.StatsContactsResponse
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/stats/contacts [get]
func (h *StatsHandler) ContactSummaries(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	rows, err := h.svc.ContactSummaries(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "stats failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rows})
}

// Summary 汇总统计（可选账本、年份；按月汇总）。
//
//	@Summary		汇总统计
//	@Tags			stats
//	@Security		BearerAuth
//	@Produce		json
//	@Param			ledger_id	query	int	false	"账本 ID"
//	@Param			year		query	int	false	"筛选年份（不传则按全部时间）"
//	@Success		200			{object}	service.SummaryResult
//	@Failure		401			{object}	handler.ErrorJSON
//	@Failure		500			{object}	handler.ErrorJSON
//	@Router			/stats/summary [get]
func (h *StatsHandler) Summary(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ledgerID := queryUintPtr(c, "ledger_id")
	year := queryIntPtr(c, "year")
	out, err := h.svc.Summary(c.Request.Context(), uid, ledgerID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "summary failed"})
		return
	}
	c.JSON(http.StatusOK, out)
}
