package services

import (
    "FINANCIALPROJECT/models"
    "gorm.io/gorm"
    "time"
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

    // Validate CreatedAt
    if moneyFlow.CreatedAt.IsZero() {
        validationErrors["created_at"] = append(validationErrors["created_at"], "created_at must be provided and valid")
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

func (s *MoneyFlowService) GetFlowsByTypeAndDate(flowType string, startDate, endDate time.Time, userID uint) ([]models.MoneyFlowDetail, error) {
    var flows []models.MoneyFlowDetail
    isIncome := flowType == "income"
    result := s.DB.Model(&models.MoneyFlow{}).
        Select("SUM(money_flows.amount) as amount, categories.name as category_name, categories.image").
        Joins("join categories on categories.id = money_flows.category_id").
        Where("money_flows.is_income = ? AND money_flows.created_at BETWEEN ? AND ? AND money_flows.user_id = ?", isIncome, startDate, endDate, userID).
        Group("categories.name, categories.image").
        Scan(&flows)

    if result.Error != nil {
        return nil, result.Error
    }
    return flows, nil
}

func (s *MoneyFlowService) GetAllFlowsByType(flowType string) ([]models.MoneyFlowDetail, error) {
    var flows []models.MoneyFlowDetail
    isIncome := flowType == "income"
    result := s.DB.Model(&models.MoneyFlow{}).
        Select("money_flows.created_at, SUM(money_flows.amount) as amount, categories.name as category_name, categories.image").
        Joins("join categories on categories.id = money_flows.category_id").
        Where("money_flows.is_income = ?", isIncome).
        Group("categories.name, categories.image, money_flows.created_at").
        Order("money_flows.created_at").
        Scan(&flows)

    if result.Error != nil {
        return nil, result.Error
    }
    return flows, nil
}

func (s *MoneyFlowService) GetUserBalance(userID uint) (float64, error) {
    var totalIncome float64
    var totalOutcome float64

    // Sumar todos los ingresos
    err := s.DB.Model(&models.MoneyFlow{}).
        Where("user_id = ? AND is_income = true", userID).
        Select("SUM(amount)").
        Row().
        Scan(&totalIncome)
    if err != nil {
        return 0, err
    }

    // Sumar todos los egresos
    err = s.DB.Model(&models.MoneyFlow{}).
        Where("user_id = ? AND is_outcome = true", userID).
        Select("SUM(amount)").
        Row().
        Scan(&totalOutcome)
    if err != nil {
        return 0, err
    }

    // Calcular el saldo final
    balance := totalIncome - totalOutcome

    return balance, nil
}