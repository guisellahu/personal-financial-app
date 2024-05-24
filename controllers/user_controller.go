package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/models"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/middleware"
    "FINANCIALPROJECT/utils"
    "github.com/dgrijalva/jwt-go"
    "strings"
)

type UserController struct {
    UserService *services.UserService
}

func NewUserController(us *services.UserService) *UserController {
    return &UserController{UserService: us}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"general": {err.Error()}})
        return
    }

    validationErrors := make(map[string][]string)

    if user.Password != user.PasswordConfirmation {
        validationErrors["password"] = append(validationErrors["password"], "passwords do not match")
    }

    serviceErrors := uc.UserService.CreateUser(&user)
    if serviceErrors != nil {
        for key, errs := range serviceErrors {
            validationErrors[key] = append(validationErrors[key], errs...)
        }
    }

    if len(validationErrors) > 0 {
        utils.SendJSONError(w, http.StatusBadRequest, validationErrors)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"general": {err.Error()}})
        return
    }

    token, err := uc.UserService.Login(creds.Email, creds.Password)
    if err != nil {
        utils.SendJSONError(w, http.StatusUnauthorized, map[string][]string{"general": {err.Error()}})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (uc *UserController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
    var req struct {
        OldPassword          string `json:"old_password"`
        NewPassword          string `json:"new_password"`
        PasswordConfirmation string `json:"password_confirmation"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"general": {err.Error()}})
        return
    }

    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        utils.SendJSONError(w, http.StatusUnauthorized, map[string][]string{"general": {"Missing Authorization header"}})
        return
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        utils.SendJSONError(w, http.StatusUnauthorized, map[string][]string{"general": {"Invalid Authorization header format"}})
        return
    }

    tokenString := parts[1]
    claims := &jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return middleware.SecretKey, nil
    })
    if err != nil || !token.Valid {
        utils.SendJSONError(w, http.StatusUnauthorized, map[string][]string{"general": {"Invalid token"}})
        return
    }

    userID := uint((*claims)["user_id"].(float64))
    validationErrors := uc.UserService.UpdatePassword(userID, req.OldPassword, req.NewPassword, req.PasswordConfirmation)
    if validationErrors != nil {
        utils.SendJSONError(w, http.StatusBadRequest, validationErrors)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}
