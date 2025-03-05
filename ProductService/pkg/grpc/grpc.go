package grpc

import (
	"ByteShop/generated/product"
	"ByteShop/generated/product_item"
	"ProductService/internal/models"
	"ProductService/pkg/database"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	product.UnimplementedProductServiceServer
	product_item.UnimplementedProductItemServiceServer
}

func StartGRPC() {
	srv := &Server{}

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}

	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(grpcServer, srv)

	log.Println("gRPC сервер запущен на :9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC: %v", err)
	}
}

func StartGRPCProductItem() {
	srv := &Server{}

	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}

	grpcServer := grpc.NewServer()
	product_item.RegisterProductItemServiceServer(grpcServer, srv)

	log.Println("gRPC сервер запущен на :9003")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC: %v", err)
	}
}

func (s *Server) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	productForDelete := models.Product{}

	err := database.DB.Model(&models.Product{}).Where("id = ?", req.ProductId).Find(&productForDelete).Error
	if err != nil {
		return &product.DeleteProductResponse{
			Status: "error",
		}, err
	}

	err = deleteFile("C://Projects/BiteShop/ProductService" + productForDelete.ImageURL)
	if err != nil {
		return &product.DeleteProductResponse{
			Status: "error",
		}, err
	}

	log.Println(req.ProductId)

	err = database.DB.Table("products").Where("id = ?", req.ProductId).Delete(&models.Product{}).Error
	if err != nil {
		log.Printf("Ошибка удаления продукта %s: %v", req.ProductId, err)
		return &product.DeleteProductResponse{
			Status: "error",
		}, err
	}
	return &product.DeleteProductResponse{
		Status: "ok",
	}, nil
}

func deleteFile(filePath string) error {

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {

		return fmt.Errorf("файл не существует: %s", filePath)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY, info.Mode())
	if err != nil {
		return fmt.Errorf("нет прав на удаление файла: %s, ошибка: %v", filePath, err)
	}
	file.Close()

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("ошибка при удалении файла: %v", err)
	}

	return nil
}

func (s *Server) DeleteAllProducts(ctx context.Context, req *product.DeleteAllProductsRequest) (*product.DeleteAllProductsResponse, error) {
	productsForDelete := []models.Product{}

	err := database.DB.Where("user_id = ?", req.UserId).Find(&productsForDelete).Error
	if err != nil {
		return &product.DeleteAllProductsResponse{
			Status: "error",
		}, err
	}

	for _, v := range productsForDelete {
		imagePath := filepath.Join("C://Projects/BiteShop/ProductService", v.ImageURL)

		deleteFile(imagePath)
	}

	err = database.DB.Where("user_id = ?", req.UserId).Delete(&models.Product{}).Error
	if err != nil {
		log.Println(err)
		return &product.DeleteAllProductsResponse{
			Status: "error",
		}, err
	}

	return &product.DeleteAllProductsResponse{
		Status: "ok",
	}, nil
}

func (s *Server) GiveProductItem(ctx context.Context, req *product_item.GiveProductItemRequest) (*product_item.GiveProductItemResponse, error) {
	productID := req.ProductID
	log.Println(productID)
	var productItem models.ProductItems
	var productData models.Product

	err := database.DB.Where("id = ?", productID).Find(&productData).Error
	if err != nil {
		return &product_item.GiveProductItemResponse{
			Status: "error",
			Value:  "",
			ID:     "",
			Name:   "",
		}, err
	}

	if productData.Type == "link" {
		return &product_item.GiveProductItemResponse{
			Status: "success",
			Value:  productData.Link,
			Name:   productData.Name,
			ID:     "",
			Type:   "link",
		}, nil
	}

	err = database.DB.Where("product_id = ? AND is_solid = ?", productID, false).First(&productItem).Error
	if err != nil {
		return &product_item.GiveProductItemResponse{
			Status: "error",
			Value:  "",
			ID:     "",
			Name:   "",
		}, err
	}

	value := productItem.Value

	updateData := map[string]interface{}{
		"is_solid": true,
		"sold_at":  time.Now(),
	}

	err = database.DB.Model(&productItem).Updates(updateData).Error
	if err != nil {
		return &product_item.GiveProductItemResponse{
			Status: "error",
			Value:  "",
			Name:   "",
			ID:     "",
		}, err
	}

	return &product_item.GiveProductItemResponse{
		Status: "success",
		Value:  value,
		Name:   productData.Name,
		ID:     productItem.ID.String(),
		Type:   "key",
	}, nil
}
