package main

import (
	"context"
	pb "demo4/proto/user"
	"demo4/server"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	conn, err := grpc.DialContext(ctx, ":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("dial error")
	}
	svc := NewClient(conn)
	output, err := svc.Get(ctx, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func NewClient(conn *grpc.ClientConn) server.UserService {
	endpoint := grpctransport.NewClient(
		conn, "user.User", "Get",
		server.EncodeGRPCGetRequest,
		server.DecodeGRPCGetResponse,
		pb.UserResponse{},
	).Endpoint()
	return server.Endpoints{GetEndpoint: endpoint}
}
