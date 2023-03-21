package main

import (
	"context"
	"fmt"
	"g09-social-todo-list/demogrpc/demo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) GetItemLikes(ctx context.Context, req *demo.GetItemLikesReq) (*demo.ItemLikesResp, error) {
	log.Println("New request with", req)
	return &demo.ItemLikesResp{Result: map[int32]int32{1: 1, 2: 4}}, nil
}

func main() {
	address := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("Server is listening on %v ...", address)

	s := grpc.NewServer()

	demo.RegisterItemLikeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
