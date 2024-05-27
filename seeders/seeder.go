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

	// Número total de registros por usuario
	totalRecords := 50

	// Ciclo para crear registros para cada usuario
	for i := 0; i < 5; i++ {
		user := models.User{
			Username: "user_" + strconv.Itoa(i),
			Email:    "user_" + strconv.Itoa(i) + "@example.com",
			Password: hashedPassword,
		}
		db.Create(&user)

		// Balance inicial para el usuario
		balance := 1000.0

		// Calcular el incremento o decremento lineal para cada transacción
		increment := balance / float64(totalRecords)

		// Ciclo para crear registros de ingresos
		for j := 0; j < totalRecords/2; j++ {
			// Calcular la cantidad y la fecha para la transacción
            amount := float64(rand.Intn(1160000) + 580000)
			createdAt := startDate.AddDate(0, 0, j)

			// Crear la transacción de ingreso
			moneyFlow := models.MoneyFlow{
				Amount:        amount,
				IsIncome:      true,
				IsOutcome:     false,
				FrequencyID:   uint(rand.Intn(5) + 1), // Suponiendo que tienes 5 frecuencias
				CategoryID:    uint(rand.Intn(2) + 7), // Sueldo o Regalos
				UserID:        user.ID,
				CreatedAt:     createdAt,
				DeactivatedAt: time.Time{},
			}
			db.Create(&moneyFlow)

			// Actualizar el balance
			balance += increment
		}

		// Balance inicial para el usuario (reset)
		balance = 1000.0

		// Ciclo para crear registros de egresos
		for j := 0; j < totalRecords/2; j++ {
			// Calcular la cantidad y la fecha para la transacción
            amount := float64(rand.Intn(1160000) + 580000)
			createdAt := startDate.AddDate(0, 0, j)

			// Crear la transacción de egreso
			moneyFlow := models.MoneyFlow{
				Amount:        amount,
				IsIncome:      false,
				IsOutcome:     true,
				FrequencyID:   uint(rand.Intn(5) + 1), // Suponiendo que tienes 5 frecuencias
				CategoryID:    uint(rand.Intn(7) + 1), // Cualquiera de las otras categorías
				UserID:        user.ID,
				CreatedAt:     createdAt,
				DeactivatedAt: time.Time{},
			}
			db.Create(&moneyFlow)

			// Actualizar el balance
			balance -= increment
        }
    }


}
