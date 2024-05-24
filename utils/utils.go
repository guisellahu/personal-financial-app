package utils

import "encoding/json"
import "net/http"

func JSONError(message string) string {
    errorResponse := map[string]string{"error": message}
    response, _ := json.Marshal(errorResponse)
    return string(response)
}

func SendJSONError(w http.ResponseWriter, status int, errors map[string][]string) {
    errorResponse := map[string]interface{}{
        "errors": errors,
    }
    response, _ := json.Marshal(errorResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(response)
}
