package handler

import (
	"net/http"

	"alioth-hrc/internal/middleware"

	"github.com/gin-gonic/gin"
)

// ReminderHandler 占位：后续可接入婚丧日期、长期未往来等提醒。
type ReminderHandler struct{}

func NewReminderHandler() *ReminderHandler {
	return &ReminderHandler{}
}

// List 提醒占位接口。
//
//	@Summary		提醒列表（占位）
//	@Tags			reminders
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	handler.ReminderListResponse
//	@Failure		401	{object}	handler.ErrorJSON
//	@Router			/reminders [get]
func (h *ReminderHandler) List(c *gin.Context) {
	if _, ok := middleware.UserID(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": []any{},
		"hint":  "reminders are not implemented yet",
	})
}
