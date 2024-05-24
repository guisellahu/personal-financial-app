package models

import (
    "gorm.io/gorm"
)

type Frequency struct {
    gorm.Model
    Name string `gorm:"not null" json:"name"`
}
