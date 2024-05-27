package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/utils"
	"github.com/dgrijalva/jwt-go"
    "time"
    "FINANCIALPROJECT/models"
    "bytes"
    "io/ioutil"
    "errors"
    "fmt"
)

type GoalController struct {
    GoalService      *services.GoalService
    MoneyFlowService *services.MoneyFlowService
}

func NewGoalController(gs *services.GoalService, mfs *services.MoneyFlowService) *GoalController {
    return &GoalController{
        GoalService:      gs,
        MoneyFlowService: mfs,
    }
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

func (gc *GoalController) GetGoalDetails(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    if err := gc.trainModel(userID); err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"train": {err.Error()}})
        return
    }

    // Obtener la fecha del último registro del usuario en money_flows
    lastFlowDate, err := gc.MoneyFlowService.GetLastFlowDate(userID)
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"last_flow_date": {err.Error()}})
        return
    }

    // Obtener la meta actual
    goal, err := gc.GoalService.GetCurrentGoal(userID)
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"goal": {"No current goal found or error retrieving goal"}})
        return
    }

    // Obtener el balance progresivo
    balanceEntries, err := gc.MoneyFlowService.GetProgressiveBalance(userID)
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"balance": {err.Error()}})
        return
    }

    // Realizar predicción
    predictionData, err := gc.PredictGoalOutcome(goal, lastFlowDate, goal.TargetDate)
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"prediction": {err.Error()}})
        return
    }

    // Debes hacer un assertion de tipo para 'predictionData["predictions"]'
    predictionsInterface, ok := predictionData["predictions"].([]map[string]interface{})
    if !ok {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"prediction": {"invalid prediction data format"}})
        return
    }

    // Extraer y reformatear predicciones
    predictions, err := formatPredictions(predictionsInterface)
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"prediction": {err.Error()}})
        return
    }

    // Formatear la respuesta
    response := map[string]interface{}{
        "goal": map[string]interface{}{
            "amount":     goal.Amount,
            "targetDate": goal.TargetDate.Format("2006-01-02"),
        },
        "balance": balanceEntries,
        "prediction": predictions,
    }
    json.NewEncoder(w).Encode(response)
}

func (gc *GoalController) trainModel(userID uint) error {
    requestData := map[string]uint{"user_id": userID}
    requestBody, err := json.Marshal(requestData)
    if err != nil {
        return err
    }

    resp, err := http.Post("http://python_app:5000/train", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to train model, received status code: %d", resp.StatusCode)
    }

    return nil
}

func (gc *GoalController) PredictGoalOutcome(goal *models.Goal, lastFlowDate, targetDate time.Time) (map[string]interface{}, error) {
    // Crear el JSON para enviar
    requestData := map[string]interface{}{
        "last_flow_date": lastFlowDate.Format("2006-01-02"),
        "target_date":    targetDate.Format("2006-01-02"),
    }
    requestBody, err := json.Marshal(requestData)
    if err != nil {
        return nil, err
    }

    // Hacer la petición al servidor de Python
    resp, err := http.Post("http://python_app:5000/predict", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var response struct {
        Predictions []map[string]interface{} `json:"predictions"`
    }
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, err
    }

    if len(response.Predictions) == 0 {
        return nil, errors.New("no predictions returned")
    }

    return map[string]interface{}{"predictions": response.Predictions}, nil
}

// Helper function to format prediction dates and ensure proper JSON output
func formatPredictions(predictions []map[string]interface{}) ([]map[string]interface{}, error) {
    var formatted []map[string]interface{}

    for _, pred := range predictions {
        dateStr, ok := pred["date"].(string)
        if !ok {
            return nil, errors.New("missing or invalid 'date' field: expected string")
        }
        parsedDate, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return nil, fmt.Errorf("invalid date format in 'date': %v", err)
        }
        predictionValue, ok := pred["prediction"]
        if !ok {
            return nil, errors.New("missing 'prediction' field")
        }
        formatted = append(formatted, map[string]interface{}{
            "date":      parsedDate.Format("2006-01-02"),
            "prediction": predictionValue,
        })
    }
    return formatted, nil
}