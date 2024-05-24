package main

import (
    "FINANCIALPROJECT/controllers"
    "FINANCIALPROJECT/models"
    "FINANCIALPROJECT/routes"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/seeders"
    "log"
    "net/http"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    dsn := "host=db user=postgres password=password dbname=mydb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database: ", err)
    }

    // Migrate the schema
    db.AutoMigrate(&models.User{}, &models.Category{}, &models.Frequency{}, &models.MoneyFlow{})

    // Seed data
    seeders.SeedData(db)

    userService := services.NewUserService(db)
    userController := controllers.NewUserController(userService)

    routes.RegisterRoutes(userController, db)

    // Serve static files from the "uploads" directory
    http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
