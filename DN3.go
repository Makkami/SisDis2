package main

import (
	"fmt"
	"log"
	"net"
	"github.com/Makkami/SisDis2/chat"
	"google.golang.org/grpc"
)

func crearServer() {
	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("Failed to listen on port 9003: %v\n", err)
	}

	s := chat.Server{}
	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("a %v", err)
	}
}

func main() {
	// Crear server
	go crearServer()

	fmt.Print("Escuchando en DN3 puerto 9003\n")
	for {

	}
}