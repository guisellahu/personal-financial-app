package controllers

import (
    "encoding/json"
    "net/http"
    "FINANCIALPROJECT/services"
    "FINANCIALPROJECT/utils"
    "time"
    "strconv"
)

type CategoryController struct {
    CategoryService *services.CategoryService
}

func NewCategoryController(cs *services.CategoryService) *CategoryController {
    return &CategoryController{CategoryService: cs}
}

func (cc *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := cc.CategoryService.GetAllCategories()
    if err != nil {
        utils.SendJSONError(w, http.StatusInternalServerError, map[string][]string{"general": {"failed to get categories"}})
        return
    }

    response := make([]map[string]interface{}, len(categories))
    for i, category := range categories {
        response[i] = map[string]interface{}{
            "id":    category.ID,
            "name":  category.Name,
            "image": category.Image,
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
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

    category, serviceErrors := cc.CategoryService.CreateCategory(name, image, imageName)
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
