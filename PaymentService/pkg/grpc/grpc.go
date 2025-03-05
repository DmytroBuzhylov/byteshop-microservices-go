package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ConnectionGRPC() *grpc.ClientConn {
	conn, err := grpc.Dial(":9002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Ошибка подключения gRPC:", err)
	}
	return conn
}
