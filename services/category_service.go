package services

import (
    "FINANCIALPROJECT/models"
    "gorm.io/gorm"
    "mime/multipart"
    "os"
    "path/filepath"
    "io"
    "errors"
)

type CategoryService struct {
    DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
    return &CategoryService{DB: db}
}

func (s *CategoryService) GetAllCategories(userID uint) ([]map[string]interface{}, error) {
    var categories []models.Category
    if err := s.DB.Where("user_id = ? OR user_id IS NULL", userID).Select("id", "name", "image", "user_id").Find(&categories).Error; err != nil {
        return nil, err
    }

    result := make([]map[string]interface{}, len(categories))
    for i, category := range categories {
        result[i] = map[string]interface{}{
            "id":         category.ID,
            "name":       category.Name,
            "image":      category.Image,
            "predefined": category.Predefined(),
        }
    }
    return result, nil
}

func (s *CategoryService) CreateCategory(name string, imageFile multipart.File, imageName string, userID uint) (*models.Category, map[string][]string) {
    validationErrors := make(map[string][]string)

    // Check for unique category name
    var existingCategory models.Category
    if err := s.DB.Where("name = ? AND (user_id = ? OR user_id IS NULL)", name, userID).First(&existingCategory).Error; err == nil {
        validationErrors["name"] = append(validationErrors["name"], "category name already exists")
    }

    if len(validationErrors) > 0 {
        return nil, validationErrors
    }

    category := models.Category{Name: name, UserID: userID}

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

func (s *CategoryService) UpdateCategory(userID uint, categoryID uint, name string, imageFile multipart.File, imageName string) (*models.Category, map[string][]string) {
    var category models.Category
    if err := s.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
        return nil, map[string][]string{"general": {"category not found or access denied"}}
    }

    // Actualizar datos
    category.Name = name
    if imageFile != nil {
        uploadPath := "uploads"
        imagePath := filepath.Join(uploadPath, imageName)
        category.Image = imagePath
        if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
            os.Mkdir(uploadPath, os.ModePerm)
        }
        outFile, err := os.Create(imagePath)
        if err != nil {
            return nil, map[string][]string{"image": {"failed to create image file: " + err.Error()}}
        }
        defer outFile.Close()
        if _, err := io.Copy(outFile, imageFile); err != nil {
            return nil, map[string][]string{"image": {"failed to copy image file: " + err.Error()}}
        }
    }

    if err := s.DB.Save(&category).Error; err != nil {
        return nil, map[string][]string{"general": {"failed to update category"}}
    }
    return &category, nil
}

func (s *CategoryService) DeleteCategory(userID uint, categoryID uint) error {
    var category models.Category
    if err := s.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
        return errors.New("category not found or access denied")
    }

    if err := s.DB.Delete(&category).Error; err != nil {
        return errors.New("failed to delete category")
    }
    return nil
}
