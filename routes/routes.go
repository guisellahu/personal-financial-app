package routes

import (
	"FINANCIALPROJECT/controllers"
	"FINANCIALPROJECT/services"
    "FINANCIALPROJECT/middleware"
    "net/http"
	"gorm.io/gorm"
)

func RegisterRoutes(uc *controllers.UserController, db *gorm.DB) {
    // Public routes
    http.HandleFunc("/api/users", uc.CreateUser)
    http.HandleFunc("/api/login", uc.Login)

    categoryService := services.NewCategoryService(db)
    categoryController := controllers.NewCategoryController(categoryService)

	frequencyService := services.NewFrequencyService(db)
    frequencyController := controllers.NewFrequencyController(frequencyService)

	moneyFlowService := services.NewMoneyFlowService(db)
    moneyFlowController := controllers.NewMoneyFlowController(moneyFlowService)

    // Protected routes
    http.Handle("/api/password", middleware.JWTMiddleware(http.HandlerFunc(uc.UpdatePassword)))

    http.Handle("/api/categories", middleware.JWTMiddleware(http.HandlerFunc(categoryController.GetAllCategories)))
    http.Handle("/api/categories/create", middleware.JWTMiddleware(http.HandlerFunc(categoryController.CreateCategory)))

	http.Handle("/api/frequencies", middleware.JWTMiddleware(http.HandlerFunc(frequencyController.GetAllFrequencies)))
	http.Handle("/api/money_flows/create", middleware.JWTMiddleware(http.HandlerFunc(moneyFlowController.CreateMoneyFlow)))
}


