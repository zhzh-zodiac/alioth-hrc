package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"

	"alioth-hrc/internal/middleware"
	"alioth-hrc/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GiftRecordHandler struct {
	svc *service.GiftRecordService
}

func NewGiftRecordHandler(svc *service.GiftRecordService) *GiftRecordHandler {
	return &GiftRecordHandler{svc: svc}
}

func parseGiftListFilter(c *gin.Context) (service.GiftRecordListFilter, error) {
	var f service.GiftRecordListFilter
	f.ContactID = queryUintPtr(c, "contact_id")
	f.LedgerID = queryUintPtr(c, "ledger_id")
	f.CategoryID = queryUintPtr(c, "category_id")
	f.Direction = c.Query("direction")
	from, err := queryDatePtr(c, "from_date")
	if err != nil {
		return f, err
	}
	to, err := queryDatePtr(c, "to_date")
	if err != nil {
		return f, err
	}
	f.From = from
	f.To = to
	return f, nil
}

// Create 新建礼金流水。
//
//	@Summary		创建礼金流水
//	@Tags			gift-records
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.GiftRecordCreate	true	"流水"
//	@Success		201		{object}	model.GiftRecord
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/gift-records [post]
func (h *GiftRecordHandler) Create(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var in service.GiftRecordCreate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.Create(c.Request.Context(), uid, in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reference"})
			return
		}
		if err.Error() == "invalid occurred_on format, use YYYY-MM-DD" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Get 礼金流水详情。
//
//	@Summary		礼金流水详情
//	@Tags			gift-records
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"流水 ID"
//	@Success		200	{object}	model.GiftRecord
//	@Failure		400	{object}	handler.ErrorJSON
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		404	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/gift-records/{id} [get]
func (h *GiftRecordHandler) Get(c *gin.Context) {
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

// List 礼金流水分页与筛选。
//
//	@Summary		礼金流水列表
//	@Tags			gift-records
//	@Security		BearerAuth
//	@Produce		json
//	@Param			contact_id	query	int		false	"联系人 ID"
//	@Param			ledger_id	query	int		false	"账本 ID"
//	@Param			category_id	query	int		false	"分类 ID"
//	@Param			direction	query	string	false	"方向: give 或 receive"
//	@Param			from_date	query	string	false	"开始日期 YYYY-MM-DD"
//	@Param			to_date		query	string	false	"结束日期 YYYY-MM-DD"
//	@Param			page		query	int		false	"页码"
//	@Param			page_size	query	int		false	"每页条数"
//	@Success		200			{object}	handler.GiftRecordListResponse
//	@Failure		400			{object}	handler.ErrorJSON
//	@Failure		401			{object}	handler.ErrorJSON
//	@Failure		500			{object}	handler.ErrorJSON
//	@Router			/gift-records [get]
func (h *GiftRecordHandler) List(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	f, err := parseGiftListFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date filter"})
		return
	}
	page := service.ParsePageQuery(c)
	list, total, err := h.svc.List(c.Request.Context(), uid, f, page)
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

// Update 更新礼金流水。
//
//	@Summary		更新礼金流水
//	@Tags			gift-records
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"流水 ID"
//	@Param			body	body		service.GiftRecordUpdate	true	"流水"
//	@Success		200		{object}	model.GiftRecord
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		404		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/gift-records/{id} [put]
func (h *GiftRecordHandler) Update(c *gin.Context) {
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
	var in service.GiftRecordUpdate
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
		if err.Error() == "invalid occurred_on format, use YYYY-MM-DD" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete 删除礼金流水。
//
//	@Summary		删除礼金流水
//	@Tags			gift-records
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	int	true	"流水 ID"
//	@Success		200	{object}	handler.OKStatus
//	@Failure		400	{object}	handler.ErrorJSON
//	@Failure		401	{object}	handler.ErrorJSON
//	@Failure		404	{object}	handler.ErrorJSON
//	@Failure		500	{object}	handler.ErrorJSON
//	@Router			/gift-records/{id} [delete]
func (h *GiftRecordHandler) Delete(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ExportCSV 导出礼金流水为 CSV（UTF-8 BOM）。
//
//	@Summary		导出礼金流水 CSV
//	@Tags			gift-records
//	@Security		BearerAuth
//	@Produce		text/csv
//	@Param			contact_id	query	int		false	"联系人 ID"
//	@Param			ledger_id	query	int		false	"账本 ID"
//	@Param			category_id	query	int		false	"分类 ID"
//	@Param			direction	query	string	false	"方向"
//	@Param			from_date	query	string	false	"开始日期 YYYY-MM-DD"
//	@Param			to_date		query	string	false	"结束日期 YYYY-MM-DD"
//	@Success		200			{string}	string	"CSV 文件"
//	@Failure		400			{object}	handler.ErrorJSON
//	@Failure		401			{object}	handler.ErrorJSON
//	@Failure		500			{object}	handler.ErrorJSON
//	@Router			/gift-records/export.csv [get]
func (h *GiftRecordHandler) ExportCSV(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	f, err := parseGiftListFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date filter"})
		return
	}
	rows, err := h.svc.ListForExport(c.Request.Context(), uid, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export failed"})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="gift_records.csv"`)
	_, _ = c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"日期", "方向", "金额(元)", "联系人", "分类", "账本", "备注"})
	for _, r := range rows {
		yuan := fmt.Sprintf("%.2f", float64(r.AmountCents)/100)
		dir := r.Direction
		if dir == "give" {
			dir = "送出"
		} else if dir == "receive" {
			dir = "收到"
		}
		_ = w.Write([]string{
			r.OccurredOn.Format("2006-01-02"),
			string(dir),
			yuan,
			r.ContactName,
			r.CategoryName,
			r.LedgerName,
			r.Note,
		})
	}
	w.Flush()
}
