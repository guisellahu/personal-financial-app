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
    var startDate, endDate time.Time
    var err error

    var flows []models.MoneyFlowDetail
    if created_at == "" {
        // Si no se proporciona created_at, busca todos los flujos de ese tipo
        flows, err = mfc.MoneyFlowService.GetAllFlowsByType(flowType)
    } else {
        dateRange := strings.Split(created_at, ",")
        if len(dateRange) != 2 {
            utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"created_at": {"Must include two dates separated by a comma [start,end]"}})
            return
        }
        // Parsear fechas
        startDate, err = time.Parse("2006-01-02", strings.TrimSpace(dateRange[0]))
        if err != nil {
            utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"created_at": {"Start date is invalid", err.Error()}})
            return
        }
        endDate, err = time.Parse("2006-01-02", strings.TrimSpace(dateRange[1]))
        if err != nil {
            utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"created_at": {"End date is invalid", err.Error()}})
            return
        }
        endDate = endDate.Add(24 * time.Hour) // Añade 24 horas para incluir todo el día

        claims := r.Context().Value("userClaims").(jwt.MapClaims)
        userID := uint(claims["user_id"].(float64))

        // Buscar los flujos por tipo y rango de fecha
        flows, err = mfc.MoneyFlowService.GetFlowsByTypeAndDate(flowType, startDate, endDate, userID)
    }

    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"general": {err.Error()}})
        return
    }

    response := map[string]interface{}{"flows": flows}
    if created_at == "" {
        // Si no hay filtro 'created_at', formatear la fecha para cada registro
        for i, flow := range flows {
            flows[i].FormattedDate = flow.CreatedAt.Format("2006-01-02") // Asignamos el string a FormattedDate
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
