package grpc

import (
	"ByteShop/generated/order"
	"ByteShop/generated/product_item"
	"OrderService/internal/models"
	"OrderService/internal/services"
	"OrderService/pkg/database"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

type Server struct {
	order.UnimplementedOrderServiceServer
	//GrpcClientProductService     product.ProductServiceClient
	GrpcClientProductItemService product_item.ProductItemServiceClient
}

//func ConnectionGRPCProductService() *grpc.ClientConn {
//	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		log.Fatalln("Ошибка подключения gRPC:", err)
//	}
//	return conn
//}

func ConnectionGRPCProductItemService() product_item.ProductItemServiceClient {
	conn, err := grpc.Dial(":9003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Ошибка подключения gRPC:", err)
	}
	return product_item.NewProductItemServiceClient(conn)
}

func StartGRPC() {
	grpcClient := ConnectionGRPCProductItemService()

	srv := &Server{
		GrpcClientProductItemService: grpcClient,
	}

	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}

	grpcServer := grpc.NewServer()
	order.RegisterOrderServiceServer(grpcServer, srv)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC: %v", err)
	}
}

func (s *Server) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	log.Println(req.ProductID)
	productResp, err := s.GrpcClientProductItemService.GiveProductItem(ctx, &product_item.GiveProductItemRequest{ProductID: req.ProductID})
	if err != nil {
		log.Printf("Ошибка при запросе к ProductService: %v", err)
		return &order.GetOrderResponse{Status: "error"}, err
	}

	email, err := services.GetUserEmail(req.BuyerID)
	if err != nil {
		return &order.GetOrderResponse{Status: "error"}, err
	}

	OrderID := uuid.New()
	BuyerID, _ := uuid.Parse(req.BuyerID)
	ProductID, _ := uuid.Parse(req.ProductID)

	orderDB := models.Order{
		ID:        OrderID,
		BuyerID:   BuyerID,
		ProductID: ProductID,
		ItemID:    productResp.ID,
		Amount:    req.Amount,
		Status:    "completed",
		CreatedAt: time.Now(),
	}

	data := services.EmailData{
		ProductName:  productResp.Name,
		ProductKey:   productResp.Value,
		OrderNumber:  OrderID.String(),
		PurchaseDate: time.Now().Format("02.01.2006 15:04"),
	}
	log.Println(productResp.Value)

	if productResp.Type == "key" {
		log.Println(1)
		err = services.SendMail(email, data, "key")
		if err != nil {
			return &order.GetOrderResponse{Status: "error"}, err
		}
	} else if productResp.Type == "link" {

		err = services.SendMail(email, data, "link")
		log.Println(err)
		if err != nil {
			return &order.GetOrderResponse{Status: "error"}, err
		}
	}

	err = database.DB.Table("orders").Create(&orderDB).Error
	if err != nil {
		return &order.GetOrderResponse{Status: "error"}, err
	}

	return &order.GetOrderResponse{Status: "OK"}, nil
}
