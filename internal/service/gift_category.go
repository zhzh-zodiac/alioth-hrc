package service

import (
	"context"
	"errors"

	"alioth-hrc/internal/model"

	"gorm.io/gorm"
)

type GiftCategoryService struct {
	db *gorm.DB
}

func NewGiftCategoryService(db *gorm.DB) *GiftCategoryService {
	return &GiftCategoryService{db: db}
}

type GiftCategoryCreate struct {
	Name string `json:"name" binding:"required,max=100"`
}

func (s *GiftCategoryService) List(ctx context.Context, userID uint) ([]model.GiftCategory, error) {
	var list []model.GiftCategory
	tx := s.db.WithContext(ctx).Where("user_id IS NULL OR user_id = ?", userID).Order("is_system DESC, id ASC")
	if err := tx.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *GiftCategoryService) Create(ctx context.Context, userID uint, in GiftCategoryCreate) (*model.GiftCategory, error) {
	uid := userID
	c := &model.GiftCategory{UserID: &uid, Name: in.Name, IsSystem: false}
	if err := s.db.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (s *GiftCategoryService) Delete(ctx context.Context, userID, id uint) error {
	var c model.GiftCategory
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&c).Error; err != nil {
		return err
	}
	if c.IsSystem {
		return errors.New("cannot delete system category")
	}
	var n int64
	if err := s.db.WithContext(ctx).Model(&model.GiftRecord{}).Where("user_id = ? AND category_id = ?", userID, id).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return errors.New("category is in use")
	}
	return s.db.WithContext(ctx).Delete(&model.GiftCategory{}, id).Error
}

func (s *GiftCategoryService) EnsureUsable(ctx context.Context, userID, categoryID uint) error {
	var c model.GiftCategory
	if err := s.db.WithContext(ctx).First(&c, categoryID).Error; err != nil {
		return err
	}
	if c.UserID == nil {
		return nil
	}
	if *c.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	return nil
}
