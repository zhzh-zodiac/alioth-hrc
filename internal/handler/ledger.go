package handler

import (
	"errors"
	"net/http"

	"alioth-hrc/internal/middleware"
	"alioth-hrc/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LedgerHandler struct {
	svc *service.LedgerService
}

func NewLedgerHandler(svc *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{svc: svc}
}

// Create 新建账本。
//
//	@Summary		创建账本
//	@Tags			ledgers
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.LedgerCreate	true	"账本"
//	@Success		201		{object}	model.Ledger
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/ledgers [post]
func (h *LedgerHandler) Create(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var in service.LedgerCreate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.Create(c.Request.Context(), uid, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Get 账本详情。
//
//	@Summary		账本详情
//	@Tags			ledgers
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"账本 ID"
//	@Success		200	{object}	model.Ledger
//	@Failure		400	{object}	handler.ErrorJSON
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		404	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/ledgers/{id} [get]
func (h *LedgerHandler) Get(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	out, err := h.svc.Get(c.Request.Context(), uid, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// List 账本列表。
//
//	@Summary		账本列表
//	@Tags			ledgers
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	handler.LedgerListResponse
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/ledgers [get]
func (h *LedgerHandler) List(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	list, err := h.svc.List(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

// Update 更新账本。
//
//	@Summary		更新账本
//	@Tags			ledgers
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"账本 ID"
//	@Param			body	body		service.LedgerUpdate	true	"账本"
//	@Success		200		{object}	model.Ledger
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		404		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/ledgers/{id} [put]
func (h *LedgerHandler) Update(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var in service.LedgerUpdate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.Update(c.Request.Context(), uid, id, in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete 删除账本（默认账本或有流水时拒绝）。
//
//	@Summary		删除账本
//	@Tags			ledgers
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"账本 ID"
//	@Success		200	{object}	handler.OKStatus
//	@Failure		400	{object}	handler.ErrorJSON
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		404	{object}	handler.ErrorJSON
//	@Failure		409	{object}	handler.ErrorJSON
//	@Router			/ledgers/{id} [delete]
func (h *LedgerHandler) Delete(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.svc.Delete(c.Request.Context(), uid, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if errors.Is(err, service.ErrCannotDeleteDefaultLedger) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
