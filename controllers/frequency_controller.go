package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/services"
	"FINANCIALPROJECT/utils"
)

type FrequencyController struct {
    FrequencyService *services.FrequencyService
}

func NewFrequencyController(fs *services.FrequencyService) *FrequencyController {
    return &FrequencyController{FrequencyService: fs}
}

func (fc *FrequencyController) GetAllFrequencies(w http.ResponseWriter, r *http.Request) {
    frequencies, err := fc.FrequencyService.GetAllFrequencies()
    if err != nil {
        http.Error(w, utils.JSONError("failed to get frequencies"), http.StatusInternalServerError)
        return
    }

    response := make([]map[string]interface{}, len(frequencies))
    for i, frequency := range frequencies {
        response[i] = map[string]interface{}{
            "id":   frequency.ID,
            "name": frequency.Name,
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
