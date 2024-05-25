package models

import (
    "gorm.io/gorm"
)

type Category struct {
    gorm.Model
    Name   string `gorm:"uniqueIndex:idx_name_user;not null" json:"name"`
    Image  string `json:"image"`
    UserID uint   `gorm:"index;default:null" json:"user_id"`
}
