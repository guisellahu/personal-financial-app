package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/models"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/utils"
    "github.com/dgrijalva/jwt-go"
    "time"
    "strings"
)

type MoneyFlowController struct {
    MoneyFlowService *services.MoneyFlowService
}

func NewMoneyFlowController(mfs *services.MoneyFlowService) *MoneyFlowController {
    return &MoneyFlowController{MoneyFlowService: mfs}
}

func (mfc *MoneyFlowController) CreateMoneyFlow(w http.ResponseWriter, r *http.Request) {
    var moneyFlow models.MoneyFlow
    if err := json.NewDecoder(r.Body).Decode(&moneyFlow); err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"general": {err.Error()}})
        return
    }

    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))
    moneyFlow.UserID = userID

    validationErrors := mfc.MoneyFlowService.CreateMoneyFlow(&moneyFlow)
    if validationErrors != nil {
        utils.SendJSONError(w, http.StatusBadRequest, validationErrors)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(moneyFlow)
}

func (mfc *MoneyFlowController) GetMoneyFlows(w http.ResponseWriter, r *http.Request) {
    flowType := r.URL.Query().Get("type")
    created_at := r.URL.Query().Get("created_at")
    dateRange := strings.Split(created_at, ",")

    if len(dateRange) != 2 {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"created_at": {"Must include two dates separated by a comma [start,end]"}})
        return
    }
    // Asumimos que las fechas están en el formato "YYYY-MM-DD"
    startDate, err := time.Parse("2006-01-02", strings.TrimSpace(dateRange[0]))
    if err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"created_at": {"Start date is invalid", err.Error()}})
        return
    }
    endDate, err := time.Parse("2006-01-02", strings.TrimSpace(dateRange[1]))
    if err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"created_at": {"End date is invalid", err.Error()}})
        return
    }
    endDate = endDate.Add(24 * time.Hour) // Añade 24 horas para incluir todo el día

    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    flows, err := mfc.MoneyFlowService.GetFlowsByTypeAndDate(flowType, startDate, endDate, userID)
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"general": {err.Error()}})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"flows": flows})
}
