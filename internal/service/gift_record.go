package service

import (
	"context"
	"errors"
	"time"

	"alioth-hrc/internal/model"

	"gorm.io/gorm"
)

type GiftRecordService struct {
	db         *gorm.DB
	contacts   *ContactService
	ledgers    *LedgerService
	categories *GiftCategoryService
}

func NewGiftRecordService(db *gorm.DB, contacts *ContactService, ledgers *LedgerService, categories *GiftCategoryService) *GiftRecordService {
	return &GiftRecordService{db: db, contacts: contacts, ledgers: ledgers, categories: categories}
}

type GiftRecordCreate struct {
	LedgerID    uint   `json:"ledger_id" binding:"required"`
	ContactID   uint   `json:"contact_id" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Direction   string `json:"direction" binding:"required,oneof=give receive"`
	AmountCents int64  `json:"amount_cents" binding:"required,min=1"`
	OccurredOn  string `json:"occurred_on" binding:"required"` // YYYY-MM-DD
	Note        string `json:"note" binding:"max=1000"`
}

type GiftRecordUpdate struct {
	LedgerID    uint   `json:"ledger_id" binding:"required"`
	ContactID   uint   `json:"contact_id" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Direction   string `json:"direction" binding:"required,oneof=give receive"`
	AmountCents int64  `json:"amount_cents" binding:"required,min=1"`
	OccurredOn  string `json:"occurred_on" binding:"required"`
	Note        string `json:"note" binding:"max=1000"`
}

type GiftRecordListFilter struct {
	ContactID  *uint
	LedgerID   *uint
	CategoryID *uint
	Direction  string
	From       *time.Time
	To         *time.Time
}

func parseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, time.Local)
}

func (s *GiftRecordService) Create(ctx context.Context, userID uint, in GiftRecordCreate) (*model.GiftRecord, error) {
	if err := s.ledgers.EnsureOwned(ctx, userID, in.LedgerID); err != nil {
		return nil, err
	}
	if err := s.contacts.EnsureOwned(ctx, userID, in.ContactID); err != nil {
		return nil, err
	}
	if err := s.categories.EnsureUsable(ctx, userID, in.CategoryID); err != nil {
		return nil, err
	}
	d, err := parseDate(in.OccurredOn)
	if err != nil {
		return nil, errors.New("invalid occurred_on format, use YYYY-MM-DD")
	}
	rec := &model.GiftRecord{
		UserID:      userID,
		LedgerID:    in.LedgerID,
		ContactID:   in.ContactID,
		CategoryID:  in.CategoryID,
		Direction:   model.GiftDirection(in.Direction),
		AmountCents: in.AmountCents,
		OccurredOn:  d,
		Note:        in.Note,
	}
	if err := s.db.WithContext(ctx).Create(rec).Error; err != nil {
		return nil, err
	}
	return rec, nil
}

func (s *GiftRecordService) Get(ctx context.Context, userID, id uint) (*model.GiftRecord, error) {
	var r model.GiftRecord
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GiftRecordService) List(ctx context.Context, userID uint, f GiftRecordListFilter, page Page) ([]model.GiftRecord, int64, error) {
	tx := s.db.WithContext(ctx).Model(&model.GiftRecord{}).Where("user_id = ?", userID)
	if f.ContactID != nil {
		tx = tx.Where("contact_id = ?", *f.ContactID)
	}
	if f.LedgerID != nil {
		tx = tx.Where("ledger_id = ?", *f.LedgerID)
	}
	if f.CategoryID != nil {
		tx = tx.Where("category_id = ?", *f.CategoryID)
	}
	if f.Direction != "" {
		tx = tx.Where("direction = ?", f.Direction)
	}
	if f.From != nil {
		tx = tx.Where("occurred_on >= ?", f.From.Format("2006-01-02"))
	}
	if f.To != nil {
		tx = tx.Where("occurred_on <= ?", f.To.Format("2006-01-02"))
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.GiftRecord
	if err := tx.Order("occurred_on DESC, id DESC").Limit(page.PageSize).Offset(page.Offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *GiftRecordService) Update(ctx context.Context, userID, id uint, in GiftRecordUpdate) (*model.GiftRecord, error) {
	var r model.GiftRecord
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&r).Error; err != nil {
		return nil, err
	}
	if err := s.ledgers.EnsureOwned(ctx, userID, in.LedgerID); err != nil {
		return nil, err
	}
	if err := s.contacts.EnsureOwned(ctx, userID, in.ContactID); err != nil {
		return nil, err
	}
	if err := s.categories.EnsureUsable(ctx, userID, in.CategoryID); err != nil {
		return nil, err
	}
	d, err := parseDate(in.OccurredOn)
	if err != nil {
		return nil, errors.New("invalid occurred_on format, use YYYY-MM-DD")
	}
	r.LedgerID = in.LedgerID
	r.ContactID = in.ContactID
	r.CategoryID = in.CategoryID
	r.Direction = model.GiftDirection(in.Direction)
	r.AmountCents = in.AmountCents
	r.OccurredOn = d
	r.Note = in.Note
	if err := s.db.WithContext(ctx).Save(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GiftRecordService) Delete(ctx context.Context, userID, id uint) error {
	res := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).Delete(&model.GiftRecord{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

const exportMaxRows = 10000

type GiftRecordExportRow struct {
	ID           uint                `json:"id" gorm:"column:id"`
	UserID       uint                `json:"user_id" gorm:"column:user_id"`
	LedgerID     uint                `json:"ledger_id" gorm:"column:ledger_id"`
	ContactID    uint                `json:"contact_id" gorm:"column:contact_id"`
	CategoryID   uint                `json:"category_id" gorm:"column:category_id"`
	Direction    model.GiftDirection `json:"direction" gorm:"column:direction"`
	AmountCents  int64               `json:"amount_cents" gorm:"column:amount_cents"`
	OccurredOn   time.Time           `json:"occurred_on" gorm:"column:occurred_on"`
	Note         string              `json:"note" gorm:"column:note"`
	CreatedAt    time.Time           `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time           `json:"updated_at" gorm:"column:updated_at"`
	ContactName  string              `json:"contact_name" gorm:"column:contact_name"`
	CategoryName string              `json:"category_name" gorm:"column:category_name"`
	LedgerName   string              `json:"ledger_name" gorm:"column:ledger_name"`
}

func (s *GiftRecordService) ListForExport(ctx context.Context, userID uint, f GiftRecordListFilter) ([]GiftRecordExportRow, error) {
	q := `
SELECT
  gr.id, gr.user_id, gr.ledger_id, gr.contact_id, gr.category_id, gr.direction, gr.amount_cents, gr.occurred_on, gr.note, gr.created_at, gr.updated_at,
  c.name AS contact_name, cat.name AS category_name, l.name AS ledger_name
FROM gift_records gr
INNER JOIN contacts c ON c.id = gr.contact_id AND c.user_id = gr.user_id
INNER JOIN gift_categories cat ON cat.id = gr.category_id
INNER JOIN ledgers l ON l.id = gr.ledger_id AND l.user_id = gr.user_id
WHERE gr.user_id = ?
`
	args := []any{userID}
	if f.ContactID != nil {
		q += " AND gr.contact_id = ?"
		args = append(args, *f.ContactID)
	}
	if f.LedgerID != nil {
		q += " AND gr.ledger_id = ?"
		args = append(args, *f.LedgerID)
	}
	if f.CategoryID != nil {
		q += " AND gr.category_id = ?"
		args = append(args, *f.CategoryID)
	}
	if f.Direction != "" {
		q += " AND gr.direction = ?"
		args = append(args, f.Direction)
	}
	if f.From != nil {
		q += " AND gr.occurred_on >= ?"
		args = append(args, f.From.Format("2006-01-02"))
	}
	if f.To != nil {
		q += " AND gr.occurred_on <= ?"
		args = append(args, f.To.Format("2006-01-02"))
	}
	q += " ORDER BY gr.occurred_on DESC, gr.id DESC LIMIT ?"
	args = append(args, exportMaxRows)

	var rows []GiftRecordExportRow
	if err := s.db.WithContext(ctx).Raw(q, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
