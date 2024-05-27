package seeders

import (
    "gorm.io/gorm"
	"FINANCIALPROJECT/models"
    "math/rand"
    "time"
    "strconv"
    "golang.org/x/crypto/bcrypt"
)

func SeedData(db *gorm.DB) {
    categories := []models.Category{
        {Name: "Transporte", Image: "uploads/transporte.png"},
        {Name: "Hogar", Image: "uploads/hogar.png"},
        {Name: "Educación", Image: "uploads/educacion.png"},
        {Name: "Comida", Image: "uploads/comida.png"},
        {Name: "Salud", Image: "uploads/salud.png"},
        {Name: "Ocio", Image: "uploads/ocio.png"},
        {Name: "Sueldo", Image: "uploads/sueldo.png"},
        {Name: "Regalos", Image: "uploads/regalos.png"},
        {Name: "Gimnasio", Image: "uploads/gimnasio.png"},
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

    // SIMULATED DATA

rand.Seed(time.Now().UnixNano())
hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("hashed_password"), bcrypt.DefaultCost)
hashedPassword := string(hashedBytes)

// Definir fechas de inicio y fin
startDate := time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC)
endDate := time.Date(2024, time.May, 26, 0, 0, 0, 0, time.UTC) // 26 para incluir el 25 de mayo

for i := 0; i < 100; i++ {
    user := models.User{
        Username: "user_" + strconv.Itoa(i),
        Email:    "user_" + strconv.Itoa(i) + "@example.com",
        Password: hashedPassword,
    }
    db.Create(&user)

    // Crear transacciones para cada usuario
    for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, rand.Intn(3)+1) { // Entre 1 y 3 días de intervalo
        createdAt := time.Date(d.Year(), d.Month(), d.Day(), rand.Intn(24), rand.Intn(60), rand.Intn(60), 0, d.Location()) // Generar un timestamp aleatorio para ese día
        amount := float64(rand.Intn(1160000) + 580000) // Montos entre 100 y 1100
        isIncome := rand.Intn(2) == 1
        moneyFlow := models.MoneyFlow{
            Amount:        amount,
            IsIncome:      isIncome,
            IsOutcome:     !isIncome,
            FrequencyID:   uint(rand.Intn(5) + 1), // Suponiendo que tienes 5 frecuencias
            CategoryID:    uint(rand.Intn(5) + 1), // Suponiendo que tienes 5 categorías
            UserID:        user.ID,
            CreatedAt:     createdAt,
            DeactivatedAt: time.Time{},
        }
        db.Create(&moneyFlow)
        }
    }


}
