package models

import (
    "gorm.io/gorm"
)

type Category struct {
    gorm.Model
    Name  string `gorm:"unique;not null" json:"name"`
    Image string `json:"image"`
}