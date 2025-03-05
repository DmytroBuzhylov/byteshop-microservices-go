package main

import (
	"AdminService/config"
	"AdminService/internal/handler"
	"AdminService/internal/middleware"
	g "AdminService/pkg/grpc"
	"ByteShop/generated/auth"
	"ByteShop/generated/product"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"log"
)

var grpcClientProductService product.ProductServiceClient
var grpcClientAuthService auth.AuthServiceClient

func main() {
	e := echo.New()
	middleware.CorsMiddleware(e)

	grpcConns := []*grpc.ClientConn{}

	grpcConnProductService := g.GRPCConnectionForProductService()
	grpcConns = append(grpcConns, grpcConnProductService)
	grpcClientProductService = product.NewProductServiceClient(grpcConnProductService)

	grpcConnAuthService := g.GRPCConnectionForAuthService()
	grpcConns = append(grpcConns, grpcConnAuthService)
	grpcClientAuthService = auth.NewAuthServiceClient(grpcConnAuthService)

	defer func() {
		for _, conn := range grpcConns {
			conn.Close()
		}
	}()

	productHandler := handler.AdminHandler{
		GrpcClientProductService: grpcClientProductService,
		GrpcClientAuthService:    grpcClientAuthService,
	}

	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.CheckAdminMiddleware)
	adminGroup.DELETE("/delete/product/:id", productHandler.RequestDeleteProductHandler)
	adminGroup.DELETE("/delete/all-products/:id", productHandler.RequestDeleteAllProductsHandler)
	adminGroup.PUT("/ban/:id", productHandler.RequestBanUserHandler)
	adminGroup.PUT("/ban/by-product-id/:id", productHandler.RequestBanUserByProductIdHandler)
	adminGroup.PUT("/:id/role", productHandler.RequestChangeRoleHandler)
	adminGroup.GET("/users", productHandler.RequestGetUsersHandler)

	log.Println(config.ServiceConfig.ADMIN_SERVICE_URL)
	if err := e.Start(config.ServiceConfig.ADMIN_SERVICE_URL); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
