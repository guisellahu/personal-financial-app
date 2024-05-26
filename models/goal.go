package models

import (
    "time"
    "gorm.io/gorm"
)

type Goal struct {
    gorm.Model
    Amount      float64   `gorm:"not null" json:"amount"`
    TargetDate  time.Time `gorm:"not null" json:"target_date"`
    UserID      uint      `gorm:"not null" json:"user_id"`
}
