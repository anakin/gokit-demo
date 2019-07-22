package main

import (
	"context"
	pb "demo4/proto/user"
	"demo4/server"
	"demo4/middleware"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"log"
	"time"
	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	commandName:=""
	hystrix.ConfigureCommand(commandName,hystrix.CommandConfig{
		Timeout: 1000*5,
		ErrorPercentThreshold: 1,
		SleepWindow: 10000,
		MaxConcurrentRequests: 1000,
		RequestVolumeThreshold: 5,
	})
	breakmw:=middleware.Hystrix(commandName)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	conn, err := grpc.DialContext(ctx, ":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("dial error")
	}
	endpoint := grpctransport.NewClient(
		conn, "user.User", "Get",
		server.EncodeGRPCGetRequest,
		server.DecodeGRPCGetResponse,
		pb.UserResponse{},
	).Endpoint()
	endpoint = breakmw(endpoint)
	svc:=server.Endpoints{GetEndpoint: endpoint}
	for i:=0;i<10;i++{
		output, err := svc.Get(ctx, 1)
			fmt.Println("当前时间: ", time.Now().Format("2006-01-02 15:04:05.99"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
	}
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
