package models

import (
    "time"
    "gorm.io/gorm"
)

type MoneyFlow struct {
    gorm.Model
    Amount        float64   `gorm:"not null" json:"amount"`
    DeactivatedAt time.Time `json:"deactivated_at,omitempty"`
    IsOutcome     bool      `json:"is_outcome"`
    IsIncome      bool      `json:"is_income"`
    FrequencyID   uint      `gorm:"not null" json:"frequency_id"`
    CategoryID    uint      `gorm:"not null" json:"category_id"`
}
