package services

import (
    "FINANCIALPROJECT/models"
    "gorm.io/gorm"
)

type MoneyFlowService struct {
    DB *gorm.DB
}

func NewMoneyFlowService(db *gorm.DB) *MoneyFlowService {
    return &MoneyFlowService{DB: db}
}

func (s *MoneyFlowService) CreateMoneyFlow(moneyFlow *models.MoneyFlow) map[string][]string {
    validationErrors := make(map[string][]string)

    // Validate Amount
    if moneyFlow.Amount <= 0 {
        validationErrors["amount"] = append(validationErrors["amount"], "amount must be greater than zero")
    }

    // Validate FrequencyID
    if err := s.DB.First(&models.Frequency{}, moneyFlow.FrequencyID).Error; err != nil {
        validationErrors["frequency_id"] = append(validationErrors["frequency_id"], "invalid frequency_id")
    }

    // Validate CategoryID
    if err := s.DB.First(&models.Category{}, moneyFlow.CategoryID).Error; err != nil {
        validationErrors["category_id"] = append(validationErrors["category_id"], "invalid category_id")
    }

    if len(validationErrors) > 0 {
        return validationErrors
    }

    // Create Money Flow
    if err := s.DB.Create(moneyFlow).Error; err != nil {
        validationErrors["general"] = append(validationErrors["general"], err.Error())
        return validationErrors
    }

    return nil
}
