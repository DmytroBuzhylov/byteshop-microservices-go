package main

import (
	"OrderService/config"
	"OrderService/internal/handlers"
	"OrderService/internal/models"
	"OrderService/pkg/database"
	"OrderService/pkg/grpc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func main() {

	err := database.DB.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/api/user-orders/:id", handlers.UserOrdersHandler)

	go grpc.StartGRPC()

	go func() {
		log.Println(config.ServiceConfig.DATABASE_URL)
	}()

	e.Start(config.ServiceConfig.ORDER_SERVICE_URL)
}
