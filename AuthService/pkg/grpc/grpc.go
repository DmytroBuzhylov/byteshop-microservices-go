package grpc

import (
	"AuthService/internal/models"
	"AuthService/pkg/database"
	"ByteShop/generated/auth"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
}

func StartGRPC() {
	srv := &Server{}

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Panicf("Ошибка запуска gRPC-сервера: %v", err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, srv)

	log.Println("gRPC сервер запущен на :9001")
	if err = grpcServer.Serve(lis); err != nil {
		log.Panicf("Ошибка работы gRPC: %v", err)
	}
}

func (s *Server) GetUser(ctx context.Context, req *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	var userList []*auth.User
	for _, u := range users {
		userList = append(userList, &auth.User{
			UserId:    u.ID.String(),
			Email:     u.Email,
			Name:      u.Name,
			Role:      u.Role,
			IsBanned:  strconv.FormatBool(u.IsBanned),
			CreatedAt: u.CreatedAt.String(),
		})
	}

	return &auth.GetUserResponse{Users: userList}, nil

}

func (s *Server) BanUser(ctx context.Context, req *auth.BanUserRequest) (*auth.BanUserResponse, error) {

	userID := req.UserId
	log.Println(userID)
	err := database.DB.Table("users").Where("id = ?", userID).UpdateColumn("is_banned", true).Error
	if err != nil {
		return &auth.BanUserResponse{
			Status: "db error",
		}, nil

	}

	return &auth.BanUserResponse{
		Status: "ok",
	}, nil
}

func (s *Server) UnBanUser(ctx context.Context, req *auth.UnBanUserRequest) (*auth.UnBanUserResponse, error) {
	userID := req.UserId

	err := database.DB.Table("users").Where("id = ?", userID).UpdateColumn("is_banned", false).Error
	if err != nil {
		return &auth.UnBanUserResponse{
			Status: "db error",
		}, nil

	}

	return &auth.UnBanUserResponse{
		Status: "ok",
	}, nil
}

func (s *Server) ChangeRole(ctx context.Context, req *auth.ChangeRoleRequest) (*auth.ChangeRoleResponse, error) {
	userId, role := req.UserId, req.Role

	err := database.DB.Table("users").Where("id = ?", userId).UpdateColumn("role", role).Error
	if err != nil {
		return &auth.ChangeRoleResponse{Status: "db error"}, nil
	}

	return &auth.ChangeRoleResponse{Status: "ok"}, nil
}
