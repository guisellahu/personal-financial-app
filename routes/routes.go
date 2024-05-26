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

    goalService := services.NewGoalService(db)
    goalController := controllers.NewGoalController(goalService)

    // Protected routes
    http.Handle("/api/password", middleware.JWTMiddleware(http.HandlerFunc(uc.UpdatePassword)))
    http.Handle("/api/user/update-username", middleware.JWTMiddleware(http.HandlerFunc(uc.UpdateUsername)))

    http.Handle("/api/categories", middleware.JWTMiddleware(http.HandlerFunc(categoryController.GetAllCategories)))
    http.Handle("/api/categories/create", middleware.JWTMiddleware(http.HandlerFunc(categoryController.CreateCategory)))
    http.Handle("/api/categories/update", middleware.JWTMiddleware(http.HandlerFunc(categoryController.UpdateCategory)))
    http.Handle("/api/categories/delete", middleware.JWTMiddleware(http.HandlerFunc(categoryController.DeleteCategory)))

	http.Handle("/api/frequencies", middleware.JWTMiddleware(http.HandlerFunc(frequencyController.GetAllFrequencies)))
	http.Handle("/api/money_flows/create", middleware.JWTMiddleware(http.HandlerFunc(moneyFlowController.CreateMoneyFlow)))
    http.Handle("/api/money_flows/balance", middleware.JWTMiddleware(http.HandlerFunc(moneyFlowController.GetBalance)))
    http.Handle("/api/money_flows", middleware.JWTMiddleware(http.HandlerFunc(moneyFlowController.GetMoneyFlows)))
    http.Handle("/api/goals/create", middleware.JWTMiddleware(http.HandlerFunc(goalController.CreateGoal)))
}


