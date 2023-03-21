package main

import (
	"context"
	"g09-social-todo-list/demogrpc/demo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	client := demo.NewItemLikeServiceClient(cc)

	for i := 1; i <= 3; i++ {
		resp, _ := client.GetItemLikes(context.Background(), &demo.GetItemLikesReq{Ids: []int32{1, 2, 3}})
		log.Println(resp.Result)
	}

}
