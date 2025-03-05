package handlers

import (
	"OrderService/internal/models"
	"OrderService/pkg/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UserOrdersHandler(c echo.Context) error {
	userID := c.Param("id")
	var userOrders []models.Order

	err := database.DB.Where("buyer_id = ?", userID).Find(&userOrders).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	return c.JSON(http.StatusOK, userOrders)
}
