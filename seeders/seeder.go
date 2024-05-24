package seeders

import (
    "gorm.io/gorm"
	"FINANCIALPROJECT/models"
)

func SeedData(db *gorm.DB) {
    categories := []models.Category{
        {Name: "Viaje"},
        {Name: "Hogar"},
        {Name: "Educación"},
        {Name: "Comida"},
        {Name: "Salud"},
    }

    for _, category := range categories {
        db.FirstOrCreate(&category, models.Category{Name: category.Name})
    }

    frequencies := []models.Frequency{
        {Name: "Esporádico"},
        {Name: "Diario"},
        {Name: "Semanal"},
        {Name: "Quincenal"},
        {Name: "Mensual"},
    }

    for _, frequency := range frequencies {
        db.FirstOrCreate(&frequency, models.Frequency{Name: frequency.Name})
    }
}
