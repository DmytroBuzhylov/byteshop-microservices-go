package main

import (
	"ProductService/config"
	"ProductService/internal/handlers"
	"ProductService/internal/middleware"
	"ProductService/internal/models"
	"ProductService/pkg/database"
	"ProductService/pkg/grpc"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	database.Connect(config.CFG.DBUrl)
	err := database.DB.AutoMigrate(
		&models.Product{},
		&models.ProductItems{},
	)
	if err != nil {
		log.Fatalln("Ошибка миграции бд")
	}

	e := echo.New()

	userHandler := handlers.UserHandler{}

	middleware.DefMiddleware(e)
	productsGroup := e.Group("/products")

	productsGroup.Use(middleware.AuthMiddleware)
	productsGroup.POST("/add-product", userHandler.CreateProduct)
	productsGroup.PUT("/product/:id", userHandler.UpdateProduct)
	e.GET("/products", userHandler.GetProducts)
	e.GET("/seller-products", userHandler.GetSellerProducts)
	e.GET("/products/:category", userHandler.GetProductsCategory)
	e.GET("/product/:id", userHandler.GetProductById)
	e.GET("/api/give-id/by-product-id/:id", userHandler.GetUserByProductId)
	e.GET("/api/give-price/by-product-id/:id", userHandler.GetPriceByProductId)
	e.GET("/api/give-all-products/:userID", userHandler.GetAllProducts)
	e.GET("/api/give-all-productID/:userID", userHandler.GetAllProductId)
	e.Static("/uploads", "uploads")

	go grpc.StartGRPC()
	go grpc.StartGRPCProductItem()

	log.Println(e.Start(config.ServiceConfig.PRODUCT_SERVICE_URL))

}
