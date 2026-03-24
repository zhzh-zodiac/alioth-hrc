package handler

import (
	"errors"
	"net/http"

	"alioth-hrc/internal/middleware"
	"alioth-hrc/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GiftCategoryHandler struct {
	svc *service.GiftCategoryService
}

func NewGiftCategoryHandler(svc *service.GiftCategoryService) *GiftCategoryHandler {
	return &GiftCategoryHandler{svc: svc}
}

// List 礼金分类（系统预设 + 当前用户自定义）。
//
//	@Summary		礼金分类列表
//	@Tags			gift-categories
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	handler.GiftCategoryListResponse
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/gift-categories [get]
func (h *GiftCategoryHandler) List(c *gin.Context) {
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

// Create 用户自定义礼金分类。
//
//	@Summary		创建礼金分类
//	@Tags			gift-categories
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.GiftCategoryCreate	true	"分类"
//	@Success		201		{object}	model.GiftCategory
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/gift-categories [post]
func (h *GiftCategoryHandler) Create(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var in service.GiftCategoryCreate
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

// Delete 删除用户自定义分类（系统分类不可删）。
//
//	@Summary		删除礼金分类
//	@Tags			gift-categories
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"分类 ID"
//	@Success		200	{object}	handler.OKStatus
//	@Failure		400	{object}	handler.ErrorJSON
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		404	{object}	handler.ErrorJSON
//	@Failure		409	{object}	handler.ErrorJSON
//	@Router			/gift-categories/{id} [delete]
func (h *GiftCategoryHandler) Delete(c *gin.Context) {
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
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
