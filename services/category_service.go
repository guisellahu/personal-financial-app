package services

import (
    "FINANCIALPROJECT/models"
    "gorm.io/gorm"
    "mime/multipart"
    "os"
    "path/filepath"
    "io"
)

type CategoryService struct {
    DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
    return &CategoryService{DB: db}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
    var categories []models.Category
    if err := s.DB.Select("id", "name", "image").Find(&categories).Error; err != nil {
        return nil, err
    }
    return categories, nil
}

func (s *CategoryService) CreateCategory(name string, imageFile multipart.File, imageName string) (*models.Category, map[string][]string) {
    validationErrors := make(map[string][]string)

    // Check for unique category name
    var existingCategory models.Category
    if err := s.DB.Where("name = ?", name).First(&existingCategory).Error; err == nil {
        validationErrors["name"] = append(validationErrors["name"], "category name already exists")
    }

    if len(validationErrors) > 0 {
        return nil, validationErrors
    }

    category := models.Category{Name: name}

    if imageFile != nil {
        uploadPath := "uploads"
        if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
            os.Mkdir(uploadPath, os.ModePerm)
        }

        imagePath := filepath.Join(uploadPath, imageName)
        category.Image = imagePath

        outFile, err := os.Create(imagePath)
        if err != nil {
            validationErrors["image"] = append(validationErrors["image"], "failed to create image file: "+err.Error())
            return nil, validationErrors
        }
        defer outFile.Close()

        if _, err := io.Copy(outFile, imageFile); err != nil {
            validationErrors["image"] = append(validationErrors["image"], "failed to copy image file: "+err.Error())
            return nil, validationErrors
        }
    }

    if err := s.DB.Create(&category).Error; err != nil {
        validationErrors["general"] = append(validationErrors["general"], "failed to save category to database: "+err.Error())
        return nil, validationErrors
    }

    return &category, nil
}
