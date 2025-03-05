package handlers

import (
	"AuthService/internal/models"
	"AuthService/internal/services"
	"AuthService/pkg/database"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
)

type UserHandler struct {
	DB *gorm.DB
}

var validate = validator.New()

func (u UserHandler) HandlerLogin(c echo.Context) error {

	var user models.User
	var req struct {
		Email    string
		Password string
	}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	var count int64
	u.DB.Table("users").Select("id, email, password_hash, role").Where("email = ?", req.Email).Count(&count).Scan(&user)

	if count == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user not registered"})
	}

	if err := services.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Incorrect password"})
	}

	token, err := services.GenerateJWT(user.ID.String(), user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})

}

func (u UserHandler) HandlerRegister(c echo.Context) error {

	var user models.User
	var req struct {
		Name     string `validate:"required,min=3"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	err = validate.Struct(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid data"})
	}

	var count int64
	u.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user already been registered"})
	}
	pass, err := services.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	user.Name = req.Name
	user.PasswordHash = pass
	user.Email = req.Email

	if err := u.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save user"})
	}

	if user.ID == uuid.Nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user ID"})
	}

	token, err := services.GenerateJWT(user.ID.String(), user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"token": token})

}

func (u *UserHandler) HandlerGetUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user ID"})
	}

	var user struct {
		Email    string
		Role     string
		IsBanned string
	}

	err := u.DB.Table("users").Select("email, role, is_banned").Where("id = ?", userID).Find(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error fetching user"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"email":     user.Email,
		"role":      user.Role,
		"is_banned": user.IsBanned,
	})
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func (u *UserHandler) ValidateJWT(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, "Missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"status": "Invalid token",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "Token is valid",
	})
}

func (u *UserHandler) CheckAdminRole(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, "Missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"status": "Invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"status": "Invalid token claims"})
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"status": "user_id not found in token"})
	}

	var user = &models.User{}

	err = u.DB.Model(&models.User{}).Where("id = ?", userID).Scan(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	if user.Role != "admin" && user.Role != "owner" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"status": "You don't have access"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "Token is valid",
	})
}

func (u *UserHandler) HandlerGiveUserData(c echo.Context) error {
	userID := services.GiveIdFromJWT(c)
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "jwt error"})
	}

	var user = &models.User{}

	err := u.DB.Model(&models.User{}).Where("id = ?", userID).Scan(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) HandlerGiveUserName(c echo.Context) error {
	userID := c.Param("userID")
	var user models.User
	err := u.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"status": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	return c.JSON(http.StatusOK, user.Name)
}

func (u *UserHandler) HandlerGiveEmail(c echo.Context) error {
	userID := c.Param("userID")
	var user models.User

	err := u.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"status": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"email": user.Email})
}

func (u *UserHandler) HandlerGiveSellerEmail(c echo.Context) error {
	sellerID := c.Param("sellerID")

	var user models.User
	err := u.DB.Table("users").Where("id = ?", sellerID).Find(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	if user.Role == "buyer" {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "user not seller"})
	}

	if user.PayPalEmail == "" {
		return c.JSON(http.StatusOK, map[string]string{"status": "user dont have PayPal email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"payPalEmail": user.PayPalEmail})
}

type AcceptTerms struct {
	AcceptTerms bool `json:"acceptTerms"`
}

func (u *UserHandler) HandlerBecomeSeller(c echo.Context) error {
	var req AcceptTerms
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "message": "invalid request"})
	}
	if req.AcceptTerms != true {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "message": "user dont accept"})
	}
	userID := services.GiveIdFromJWT(c)
	if userID == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}
	err := database.DB.Table("users").Where("id = ?", userID).UpdateColumn("role", "seller").Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	return c.JSON(http.StatusCreated, map[string]bool{"success": true})
}

func (u *UserHandler) HandletUserExists(c echo.Context) error {
	userID := c.Param("id")
	log.Println(userID)
	var user models.User

	err := u.DB.Table("users").Where("id = ?", userID).Find(&user).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "message": "user dont accept"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})

}

func (u *UserHandler) HandlerGiveUserId(c echo.Context) error {
	userID := services.GiveIdFromJWT(c)

	return c.JSON(http.StatusOK, map[string]string{"userID": userID})
}
