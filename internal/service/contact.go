package service

import (
	"context"
	"errors"

	"alioth-hrc/internal/model"

	"gorm.io/gorm"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(db *gorm.DB) *ContactService {
	return &ContactService{db: db}
}

type ContactCreate struct {
	Name         string `json:"name" binding:"required,max=200"`
	RelationNote string `json:"relation_note" binding:"max=200"`
	Remark       string `json:"remark" binding:"max=500"`
}

type ContactUpdate struct {
	Name         string `json:"name" binding:"required,max=200"`
	RelationNote string `json:"relation_note" binding:"max=200"`
	Remark       string `json:"remark" binding:"max=500"`
}

func (s *ContactService) Create(ctx context.Context, userID uint, in ContactCreate) (*model.Contact, error) {
	c := &model.Contact{
		UserID:       userID,
		Name:         in.Name,
		RelationNote: in.RelationNote,
		Remark:       in.Remark,
	}
	if err := s.db.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ContactService) Get(ctx context.Context, userID, id uint) (*model.Contact, error) {
	var c model.Contact
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *ContactService) List(ctx context.Context, userID uint, q string, page Page) ([]model.Contact, int64, error) {
	tx := s.db.WithContext(ctx).Model(&model.Contact{}).Where("user_id = ?", userID)
	if q != "" {
		like := "%" + q + "%"
		tx = tx.Where("name LIKE ? OR relation_note LIKE ?", like, like)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Contact
	if err := tx.Order("id DESC").Limit(page.PageSize).Offset(page.Offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *ContactService) Update(ctx context.Context, userID, id uint, in ContactUpdate) (*model.Contact, error) {
	var c model.Contact
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&c).Error; err != nil {
		return nil, err
	}
	c.Name = in.Name
	c.RelationNote = in.RelationNote
	c.Remark = in.Remark
	if err := s.db.WithContext(ctx).Save(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

var ErrContactHasRecords = errors.New("contact has gift records; delete or reassign records first")

func (s *ContactService) Delete(ctx context.Context, userID, id uint) error {
	var n int64
	if err := s.db.WithContext(ctx).Model(&model.GiftRecord{}).Where("user_id = ? AND contact_id = ?", userID, id).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return ErrContactHasRecords
	}
	res := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).Delete(&model.Contact{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *ContactService) EnsureOwned(ctx context.Context, userID, contactID uint) error {
	var n int64
	if err := s.db.WithContext(ctx).Model(&model.Contact{}).Where("user_id = ? AND id = ?", userID, contactID).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
