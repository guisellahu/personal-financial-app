package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/models"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/utils"
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

    validationErrors := mfc.MoneyFlowService.CreateMoneyFlow(&moneyFlow)
    if validationErrors != nil {
        utils.SendJSONError(w, http.StatusBadRequest, validationErrors)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(moneyFlow)
}
