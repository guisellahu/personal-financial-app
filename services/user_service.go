package services

import (
    "FINANCIALPROJECT/models"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "regexp"
    "github.com/dgrijalva/jwt-go"
    "time"
)

type UserService struct {
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *models.User) map[string][]string {
    validationErrors := make(map[string][]string)

    // Validate username
    matched, _ := regexp.MatchString(`^\S+$`, user.Username)
    if !matched {
        validationErrors["username"] = append(validationErrors["username"], "username cannot contain spaces")
    }

    // Validate email format
    matched, _ = regexp.MatchString(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, user.Email)
    if !matched {
        validationErrors["email"] = append(validationErrors["email"], "invalid email format")
    }

    // Check for unique username
    var count int64
    s.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
    if count > 0 {
        validationErrors["username"] = append(validationErrors["username"], "username already exists")
    }

    // Check for unique email
    s.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
    if count > 0 {
        validationErrors["email"] = append(validationErrors["email"], "email already exists")
    }

    if len(validationErrors) > 0 {
        return validationErrors
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        validationErrors["password"] = append(validationErrors["password"], "failed to hash password")
        return validationErrors
    }
    user.Password = string(hashedPassword)

    // Create user
    if err := s.DB.Create(user).Error; err != nil {
        validationErrors["general"] = append(validationErrors["general"], err.Error())
        return validationErrors
    }

    return nil
}

func (s *UserService) Login(email, password string) (string, error) {
    var user models.User
    if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return "", errors.New("invalid email or password")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid email or password")
    }

    // Generate JWT token
    token, err := generateJWT(user)
    if err != nil {
        return "", err
    }

    return token, nil
}

func (s *UserService) UpdatePassword(userID uint, oldPassword, newPassword, passwordConfirmation string) map[string][]string {
    var user models.User
    validationErrors := make(map[string][]string)

    if err := s.DB.First(&user, userID).Error; err != nil {
        validationErrors["general"] = append(validationErrors["general"], "user not found")
        return validationErrors
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
        validationErrors["old_password"] = append(validationErrors["old_password"], "incorrect old password")
        return validationErrors
    }

    if newPassword != passwordConfirmation {
        validationErrors["new_password"] = append(validationErrors["new_password"], "new password and confirmation do not match")
        return validationErrors
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        validationErrors["new_password"] = append(validationErrors["new_password"], "failed to hash new password")
        return validationErrors
    }

    user.Password = string(hashedPassword)

    if err := s.DB.Save(&user).Error; err != nil {
        validationErrors["general"] = append(validationErrors["general"], "failed to update password")
        return validationErrors
    }

    return nil
}

func generateJWT(user models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })

    secretKey := "your_secret_key" // Replace this with your secret key
    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *UserService) UpdateUsername(userID uint, newUsername string) map[string][]string {
    validationErrors := make(map[string][]string)

    // Validate new username format
    matched, _ := regexp.MatchString(`^\S+$`, newUsername)
    if !matched {
        validationErrors["username"] = []string{"username cannot contain spaces"}
        return validationErrors
    }

    // Check for unique username
    var count int64
    s.DB.Model(&models.User{}).Where("username = ?", newUsername).Not("id = ?", userID).Count(&count)
    if count > 0 {
        validationErrors["username"] = []string{"username already exists"}
        return validationErrors
    }

    // Update username
    var user models.User
    if err := s.DB.First(&user, userID).Error; err != nil {
        validationErrors["general"] = []string{"user not found"}
        return validationErrors
    }

    user.Username = newUsername
    if err := s.DB.Save(&user).Error; err != nil {
        validationErrors["general"] = []string{"failed to update username"}
        return validationErrors
    }

    return nil
}
