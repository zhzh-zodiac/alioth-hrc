package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StatsService struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewStatsService(db *gorm.DB, rdb *redis.Client) *StatsService {
	return &StatsService{db: db, rdb: rdb}
}

type ContactStatRow struct {
	ContactID      uint       `json:"contact_id"`
	TotalReceive   int64      `json:"total_receive_cents"`
	TotalGive      int64      `json:"total_give_cents"`
	BalanceCents   int64      `json:"balance_cents"`
	LastOccurredOn *time.Time `json:"last_occurred_on,omitempty"`
}

func (s *StatsService) ContactSummaries(ctx context.Context, userID uint) ([]ContactStatRow, error) {
	type row struct {
		ContactID    uint
		TotalReceive int64
		TotalGive    int64
		LastDate     *time.Time `gorm:"column:last_date"`
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
SELECT
  contact_id AS contact_id,
  COALESCE(SUM(CASE WHEN direction = 'receive' THEN amount_cents ELSE 0 END), 0) AS total_receive,
  COALESCE(SUM(CASE WHEN direction = 'give' THEN amount_cents ELSE 0 END), 0) AS total_give,
  MAX(occurred_on) AS last_date
FROM gift_records
WHERE user_id = ?
GROUP BY contact_id
ORDER BY contact_id
`, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]ContactStatRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, ContactStatRow{
			ContactID:      r.ContactID,
			TotalReceive:   r.TotalReceive,
			TotalGive:      r.TotalGive,
			BalanceCents:   r.TotalReceive - r.TotalGive,
			LastOccurredOn: r.LastDate,
		})
	}
	return out, nil
}

type MonthlyStatRow struct {
	YearMonth    string `json:"year_month"`
	ReceiveCents int64  `json:"receive_cents"`
	GiveCents    int64  `json:"give_cents"`
}

type SummaryResult struct {
	LedgerID     *uint            `json:"ledger_id,omitempty"`
	Year         *int             `json:"year,omitempty"`
	Monthly      []MonthlyStatRow `json:"monthly"`
	TotalReceive int64            `json:"total_receive_cents"`
	TotalGive    int64            `json:"total_give_cents"`
	BalanceCents int64            `json:"balance_cents"`
}

func summaryCacheKey(userID uint, ledgerID *uint, year *int) string {
	h := sha256.New()
	fmt.Fprintf(h, "u=%d", userID)
	if ledgerID != nil {
		fmt.Fprintf(h, "&l=%d", *ledgerID)
	} else {
		h.Write([]byte("&l=all"))
	}
	if year != nil {
		fmt.Fprintf(h, "&y=%d", *year)
	} else {
		h.Write([]byte("&y=all"))
	}
	return "hrc:stats:" + hex.EncodeToString(h.Sum(nil))[:32]
}

func (s *StatsService) Summary(ctx context.Context, userID uint, ledgerID *uint, year *int) (*SummaryResult, error) {
	key := summaryCacheKey(userID, ledgerID, year)
	if raw, err := s.rdb.Get(ctx, key).Bytes(); err == nil && len(raw) > 0 {
		var cached SummaryResult
		if err := json.Unmarshal(raw, &cached); err == nil {
			return &cached, nil
		}
	}

	var args []any
	args = append(args, userID)
	q := `
SELECT
  COALESCE(SUM(CASE WHEN direction = 'receive' THEN amount_cents ELSE 0 END), 0) AS receive,
  COALESCE(SUM(CASE WHEN direction = 'give' THEN amount_cents ELSE 0 END), 0) AS give
FROM gift_records
WHERE user_id = ?
`
	if ledgerID != nil {
		q += " AND ledger_id = ?"
		args = append(args, *ledgerID)
	}
	if year != nil {
		start := time.Date(*year, 1, 1, 0, 0, 0, 0, time.Local)
		end := time.Date(*year+1, 1, 1, 0, 0, 0, 0, time.Local)
		q += " AND occurred_on >= ? AND occurred_on < ?"
		args = append(args, start.Format("2006-01-02"), end.Format("2006-01-02"))
	}

	type agg struct {
		Receive int64 `gorm:"column:receive"`
		Give    int64 `gorm:"column:give"`
	}
	var total agg
	if err := s.db.WithContext(ctx).Raw(q, args...).Scan(&total).Error; err != nil {
		return nil, err
	}

	margs := []any{userID}
	mq := `
SELECT
  DATE_FORMAT(occurred_on, '%Y-%m') AS ym,
  COALESCE(SUM(CASE WHEN direction = 'receive' THEN amount_cents ELSE 0 END), 0) AS receive,
  COALESCE(SUM(CASE WHEN direction = 'give' THEN amount_cents ELSE 0 END), 0) AS give
FROM gift_records
WHERE user_id = ?
`
	if ledgerID != nil {
		mq += " AND ledger_id = ?"
		margs = append(margs, *ledgerID)
	}
	if year != nil {
		start := time.Date(*year, 1, 1, 0, 0, 0, 0, time.Local)
		end := time.Date(*year+1, 1, 1, 0, 0, 0, 0, time.Local)
		mq += " AND occurred_on >= ? AND occurred_on < ?"
		margs = append(margs, start.Format("2006-01-02"), end.Format("2006-01-02"))
	}
	mq += " GROUP BY DATE_FORMAT(occurred_on, '%Y-%m') ORDER BY ym ASC"

	type mrow struct {
		YM      string `gorm:"column:ym"`
		Receive int64  `gorm:"column:receive"`
		Give    int64  `gorm:"column:give"`
	}
	var monthly []mrow
	if err := s.db.WithContext(ctx).Raw(mq, margs...).Scan(&monthly).Error; err != nil {
		return nil, err
	}

	monthlyOut := make([]MonthlyStatRow, 0, len(monthly))
	for _, m := range monthly {
		monthlyOut = append(monthlyOut, MonthlyStatRow{
			YearMonth:    m.YM,
			ReceiveCents: m.Receive,
			GiveCents:    m.Give,
		})
	}

	res := &SummaryResult{
		LedgerID:     ledgerID,
		Year:         year,
		Monthly:      monthlyOut,
		TotalReceive: total.Receive,
		TotalGive:    total.Give,
		BalanceCents: total.Receive - total.Give,
	}

	if payload, err := json.Marshal(res); err == nil {
		_ = s.rdb.Set(ctx, key, payload, 60*time.Second).Err()
	}
	return res, nil
}
