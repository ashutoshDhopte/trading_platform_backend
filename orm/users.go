package orm

import (
	"time"
)

type Users struct {
	UserID           int64 `gorm:"primaryKey"`
	Username         string
	Email            string
	HashedPassword   string
	Salt             string
	CashBalanceCents int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	NotificationsOn  bool
}
