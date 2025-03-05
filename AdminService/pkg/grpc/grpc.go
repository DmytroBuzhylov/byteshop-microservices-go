package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GRPCConnectionForProductService() *grpc.ClientConn {
	var err error
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC-серверу: %v", err)
	}

	return conn

}

func GRPCConnectionForAuthService() *grpc.ClientConn {
	var err error
	conn, err := grpc.Dial(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC-серверу: %v", err)
	}

	return conn

}
