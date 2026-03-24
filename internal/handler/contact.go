package handler

import (
	"errors"
	"net/http"
	"strconv"

	"alioth-hrc/internal/middleware"
	"alioth-hrc/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContactHandler struct {
	svc *service.ContactService
}

func NewContactHandler(svc *service.ContactService) *ContactHandler {
	return &ContactHandler{svc: svc}
}

// Create 新建联系人。
//
//	@Summary		创建联系人
//	@Tags			contacts
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.ContactCreate	true	"联系人"
//	@Success		201		{object}	model.Contact
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/contacts [post]
func (h *ContactHandler) Create(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var in service.ContactCreate
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

// Get 联系人详情。
//
//	@Summary		联系人详情
//	@Tags			contacts
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"联系人 ID"
//	@Success		200	{object}	model.Contact
//	@Failure		400	{object}	handler.ErrorJSON
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		404	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/contacts/{id} [get]
func (h *ContactHandler) Get(c *gin.Context) {
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

// List 联系人分页列表。
//
//	@Summary		联系人列表
//	@Tags			contacts
//	@Security		BearerAuth
//	@Produce		json
//	@Param			q			query		string	false	"按姓名/关系模糊搜索"
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页条数"
//	@Success		200			{object}	handler.ContactListResponse
//	@Failure		401			{object}	handler.ErrorJSON
//	@Failure		500			{object}	handler.ErrorJSON
//	@Router			/contacts [get]
func (h *ContactHandler) List(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	q := c.Query("q")
	page := service.ParsePageQuery(c)
	list, total, err := h.svc.List(c.Request.Context(), uid, q, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":     list,
		"total":     total,
		"page":      page.Page,
		"page_size": page.PageSize,
	})
}

// Update 更新联系人。
//
//	@Summary		更新联系人
//	@Tags			contacts
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"联系人 ID"
//	@Param			body	body		service.ContactUpdate	true	"联系人"
//	@Success		200		{object}	model.Contact
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		404		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/contacts/{id} [put]
func (h *ContactHandler) Update(c *gin.Context) {
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
	var in service.ContactUpdate
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

// Delete 删除联系人（若存在礼金流水则拒绝）。
//
//	@Summary		删除联系人
//	@Tags			contacts
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"联系人 ID"
//	@Success		200	{object}	handler.OKStatus
//	@Failure		400	{object}	handler.ErrorJSON
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		404	{object}	handler.ErrorJSON
//	@Failure		409	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/contacts/{id} [delete]
func (h *ContactHandler) Delete(c *gin.Context) {
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
		if errors.Is(err, service.ErrContactHasRecords) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func parseUintParam(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
