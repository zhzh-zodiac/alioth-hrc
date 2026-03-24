package service

import (
	"context"
	"errors"

	"alioth-hrc/internal/model"

	"gorm.io/gorm"
)

type LedgerService struct {
	db *gorm.DB
}

func NewLedgerService(db *gorm.DB) *LedgerService {
	return &LedgerService{db: db}
}

type LedgerCreate struct {
	Name string `json:"name" binding:"required,max=200"`
}

type LedgerUpdate struct {
	Name string `json:"name" binding:"required,max=200"`
}

func (s *LedgerService) Create(ctx context.Context, userID uint, in LedgerCreate) (*model.Ledger, error) {
	l := &model.Ledger{UserID: userID, Name: in.Name, IsDefault: false}
	if err := s.db.WithContext(ctx).Create(l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func (s *LedgerService) Get(ctx context.Context, userID, id uint) (*model.Ledger, error) {
	var l model.Ledger
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *LedgerService) List(ctx context.Context, userID uint) ([]model.Ledger, error) {
	var list []model.Ledger
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *LedgerService) Update(ctx context.Context, userID, id uint, in LedgerUpdate) (*model.Ledger, error) {
	var l model.Ledger
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&l).Error; err != nil {
		return nil, err
	}
	l.Name = in.Name
	if err := s.db.WithContext(ctx).Save(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

var ErrCannotDeleteDefaultLedger = errors.New("cannot delete default ledger")

func (s *LedgerService) Delete(ctx context.Context, userID, id uint) error {
	var l model.Ledger
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&l).Error; err != nil {
		return err
	}
	if l.IsDefault {
		return ErrCannotDeleteDefaultLedger
	}
	var n int64
	if err := s.db.WithContext(ctx).Model(&model.GiftRecord{}).Where("user_id = ? AND ledger_id = ?", userID, id).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return errors.New("ledger has gift records")
	}
	return s.db.WithContext(ctx).Delete(&model.Ledger{}, id).Error
}

func (s *LedgerService) EnsureOwned(ctx context.Context, userID, ledgerID uint) error {
	var n int64
	if err := s.db.WithContext(ctx).Model(&model.Ledger{}).Where("user_id = ? AND id = ?", userID, ledgerID).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
