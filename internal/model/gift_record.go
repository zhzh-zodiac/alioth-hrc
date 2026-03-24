package model

import (
	"time"
)

type GiftDirection string

const (
	DirectionGive    GiftDirection = "give"
	DirectionReceive GiftDirection = "receive"
)

type GiftRecord struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	UserID      uint          `gorm:"index;not null" json:"user_id"`
	LedgerID    uint          `gorm:"index;not null" json:"ledger_id"`
	ContactID   uint          `gorm:"index;not null" json:"contact_id"`
	CategoryID  uint          `gorm:"index;not null" json:"category_id"`
	Direction   GiftDirection `gorm:"size:16;not null" json:"direction"`
	AmountCents int64         `gorm:"not null" json:"amount_cents"`
	OccurredOn  time.Time     `gorm:"type:date;not null" json:"occurred_on"`
	Note        string        `gorm:"size:1000" json:"note"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
