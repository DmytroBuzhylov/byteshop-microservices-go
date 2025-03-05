package main

import (
	"ByteShop/generated/order"
	"PaymentService/config"
	"PaymentService/internal/handlers"
	"PaymentService/internal/models"
	"PaymentService/pkg/database"
	g "PaymentService/pkg/grpc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

const (
	platformFee = 0.1
)

var paypalAPI = config.ServiceConfig.PAYPAL_API

var grpcClientOrderService order.OrderServiceClient

func main() {
	grpcConnOrderService := g.ConnectionGRPC()
	grpcClientOrderService = order.NewOrderServiceClient(grpcConnOrderService)
	defer grpcConnOrderService.Close()
	database.Connect(config.ServiceConfig.DATABASE_URL)
	e := echo.New()
	err := database.DB.AutoMigrate(&models.Payment{})
	if err != nil {
		log.Fatalln("Ошибка миграции")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	paymentHandler := handlers.PaymentHandler{
		GrpcClientPaymentService: grpcClientOrderService,
	}

	e.POST("/paypal/create-order", paymentHandler.CreatePayment)
	e.POST("/paypal/capture-order", paymentHandler.CapturePayment)

	log.Println(e.Start(config.ServiceConfig.PAYMENT_SERVICE_URL))
}
