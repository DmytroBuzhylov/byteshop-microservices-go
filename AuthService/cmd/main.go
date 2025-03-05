package main

import (
	"AuthService/config"
	"AuthService/internal/handlers"
	"AuthService/internal/middlewares"
	"AuthService/internal/models"
	"AuthService/pkg/database"
	"AuthService/pkg/grpc"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	database.Connect(cfg.DBUrl)
	e := echo.New()

	DB := database.GiveDB()
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Migrate error: %d", err)
	}

	userHandler := handlers.UserHandler{DB: DB}

	middlewares.DefMiddleware(e)

	e.POST("/auth/register", userHandler.HandlerRegister)
	e.POST("/auth/login", userHandler.HandlerLogin)
	e.GET("/auth/validate", userHandler.ValidateJWT)
	e.GET("/auth/admin", userHandler.CheckAdminRole)
	e.GET("/api/me", userHandler.HandlerGiveUserData)
	e.GET("/api/name/:userID", userHandler.HandlerGiveUserName)
	e.GET("/api/email/:userID", userHandler.HandlerGiveEmail)
	e.GET("/api/user/:sellerID/paypal", userHandler.HandlerGiveSellerEmail)
	e.GET("/api/user-exists/:id", userHandler.HandletUserExists)
	e.GET("/api/give-user-id", userHandler.HandlerGiveUserId)
	e.POST("/api/user/become-seller", userHandler.HandlerBecomeSeller)

	authGroup := e.Group("/user")
	authGroup.Use(middlewares.JWTMiddleware)
	authGroup.GET("/profile", userHandler.HandlerGetUser)
	authGroup.POST("/addProduct", userHandler.HandlerGetUser)

	go grpc.StartGRPC()

	log.Println(e.Start(cfg.AUTH_SERVICE_URL))
}
