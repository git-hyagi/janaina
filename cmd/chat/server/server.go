package main

import (
	"github.com/git-hyagi/janaina/pkg/chat"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("Failed to bind port 9001: %v", err)
	}

	chatService := chat.Server{}
	grpcServer := grpc.NewServer()
	chat.RegisterChatServer(grpcServer, &chatService)
	log.Printf("connections: %v", chatService)
	grpcServer.Serve(listener)

}
