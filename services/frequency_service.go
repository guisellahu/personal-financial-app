package services

import (
    "FINANCIALPROJECT/models"
    "gorm.io/gorm"
)

type FrequencyService struct {
    DB *gorm.DB
}

func NewFrequencyService(db *gorm.DB) *FrequencyService {
    return &FrequencyService{DB: db}
}

func (s *FrequencyService) GetAllFrequencies() ([]models.Frequency, error) {
    var frequencies []models.Frequency
    if err := s.DB.Select("id", "name").Find(&frequencies).Error; err != nil {
        return nil, err
    }
    return frequencies, nil
}
