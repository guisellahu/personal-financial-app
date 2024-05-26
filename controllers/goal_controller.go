package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/utils"
	"github.com/dgrijalva/jwt-go"
    "time"
)

type GoalController struct {
    GoalService *services.GoalService
}

func NewGoalController(gs *services.GoalService) *GoalController {
    return &GoalController{GoalService: gs}
}

func (gc *GoalController) CreateGoal(w http.ResponseWriter, r *http.Request) {
    var data struct {
        Amount     float64 `json:"amount"`
        TargetDate string  `json:"target_date"`
    }
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"general": {err.Error()}})
        return
    }

    // Parsear la fecha
    parsedDate, err := time.Parse("2006-01-02", data.TargetDate)
    if err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"target_date": {"invalid date format"}})
        return
    }

	parsedDate = parsedDate.Add(24 * time.Hour)

    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    if err := gc.GoalService.CreateGoal(userID, data.Amount, parsedDate); err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"general": {err.Error()}})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Goal created successfully"})
}
