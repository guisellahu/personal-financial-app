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
