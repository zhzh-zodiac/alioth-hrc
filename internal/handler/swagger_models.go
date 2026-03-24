package handler

import (
	"alioth-hrc/internal/model"
	"alioth-hrc/internal/service"
)

// ErrorJSON 通用错误响应（文档用）。
type ErrorJSON struct {
	Error string `json:"error"`
}

// OKStatus 通用成功状态。
type OKStatus struct {
	Status string `json:"status"`
}

// RegisterResponse 注册成功响应。
type RegisterResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ContactListResponse 联系人分页。
type ContactListResponse struct {
	Items    []model.Contact `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// LedgerListResponse 账本列表。
type LedgerListResponse struct {
	Items []model.Ledger `json:"items"`
}

// GiftCategoryListResponse 礼金分类列表。
type GiftCategoryListResponse struct {
	Items []model.GiftCategory `json:"items"`
}

// GiftRecordListResponse 礼金流水分页。
type GiftRecordListResponse struct {
	Items    []model.GiftRecord `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// StatsContactsResponse 按联系人统计。
type StatsContactsResponse struct {
	Items []service.ContactStatRow `json:"items"`
}

// ReminderListResponse 提醒占位响应。
type ReminderListResponse struct {
	Items []interface{} `json:"items"`
	Hint  string        `json:"hint"`
}
