package services

import (
    "FINANCIALPROJECT/models"
    "gorm.io/gorm"
    "time"
    "errors"
)

type GoalService struct {
    DB *gorm.DB
}

func NewGoalService(db *gorm.DB) *GoalService {
    return &GoalService{DB: db}
}

func (s *GoalService) CreateGoal(userID uint, amount float64, targetDate time.Time) error {
    var existingGoals int64
    s.DB.Model(&models.Goal{}).Where("user_id = ? AND target_date > ?", userID, time.Now()).Count(&existingGoals)
    if existingGoals > 0 {
        return errors.New("there is already an active goal with a future target date")
    }

    goal := models.Goal{
        UserID: userID,
        Amount: amount,
        TargetDate: targetDate,
    }

    if err := s.DB.Create(&goal).Error; err != nil {
        return err
    }

    return nil
}

func (s *GoalService) GetCurrentGoal(userID uint) (*models.Goal, error) {
    var goal models.Goal
    err := s.DB.Where("user_id = ? AND target_date > ?", userID, time.Now()).Order("target_date ASC").First(&goal).Error
    if err != nil {
        return nil, err
    }
    return &goal, nil
}

func (s *MoneyFlowService) GetProgressiveBalance(userID uint) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    var balance float64 = 0
    var flows []models.MoneyFlow

    // Obtener todos los flujos de dinero ordenados por fecha
    result := s.DB.Where("user_id = ?", userID).Order("created_at asc").Find(&flows)
    if result.Error != nil {
        return nil, result.Error
    }

    // Calcular el balance progresivo
    for _, flow := range flows {
        if flow.IsIncome {
            balance += flow.Amount
        } else {
            balance -= flow.Amount
        }
        results = append(results, map[string]interface{}{
            "balance":   balance,
            "createdAt": flow.CreatedAt.Format("2006-01-02"),
        })
    }

    return results, nil
}