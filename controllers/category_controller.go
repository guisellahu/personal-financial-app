package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/utils"
    "time"
    "strconv"
    "github.com/dgrijalva/jwt-go"
)

type CategoryController struct {
    CategoryService *services.CategoryService
}

func NewCategoryController(cs *services.CategoryService) *CategoryController {
    return &CategoryController{CategoryService: cs}
}

func (cc *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))
    categories, err := cc.CategoryService.GetAllCategories(userID)
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"general": {"failed to get categories"}})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"categories": categories}) // Correct usage
}


func (cc *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("name")
    image, imageHeader, err := r.FormFile("image")

    validationErrors := make(map[string][]string)

    if name == "" {
        validationErrors["name"] = append(validationErrors["name"], "name is required")
    }

    if err != nil && err != http.ErrMissingFile {
        validationErrors["image"] = append(validationErrors["image"], "failed to process image file")
    }

    var imageName string
    if imageHeader != nil {
        imageName = strconv.FormatInt(time.Now().Unix(), 10) + "_" + imageHeader.Filename
    }

    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    category, serviceErrors := cc.CategoryService.CreateCategory(name, image, imageName, userID)
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
    json.NewEncoder(w).Encode(category)
}

func (cc *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
    categoryID, _ := strconv.Atoi(r.URL.Query().Get("id"))
    name := r.FormValue("name")
    image, imageHeader, _ := r.FormFile("image")

    var imageName string
    if imageHeader != nil {
        imageName = strconv.FormatInt(time.Now().Unix(), 10) + "_" + imageHeader.Filename
    }

    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    category, err := cc.CategoryService.UpdateCategory(userID, uint(categoryID), name, image, imageName)
    if err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(category)
}

func (cc *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
    categoryID, _ := strconv.Atoi(r.URL.Query().Get("id"))

    claims := r.Context().Value("userClaims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    if err := cc.CategoryService.DeleteCategory(userID, uint(categoryID)); err != nil {
        utils.SendJSONError(w, http.StatusBadRequest, map[string][]string{"general": {err.Error()}})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
}