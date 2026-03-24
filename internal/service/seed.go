package service

import (
	"alioth-hrc/internal/model"

	"gorm.io/gorm"
)

var defaultCategoryNames = []string{"婚礼", "白事", "生日", "满月", "其他"}

func SeedGiftCategories(db *gorm.DB) error {
	for _, name := range defaultCategoryNames {
		var n int64
		if err := db.Model(&model.GiftCategory{}).
			Where("user_id IS NULL AND name = ?", name).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			continue
		}
		c := model.GiftCategory{Name: name, IsSystem: true}
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}
	return nil
}
