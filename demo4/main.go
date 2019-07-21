package main

import (
	pb "demo4/proto/user"
	"demo4/server"
	"fmt"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	svc := server.New()
	endpoint := server.Endpoints{
		GetEndpoint: server.MakeUserEndpoint(svc),
	}
	listener,err:=net.Listen("tcp",":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ser:=server.MakeGRPCServer(endpoint)
	s:=grpc.NewServer()
	pb.RegisterUserServer(s,ser)
	err=s.Serve(listener)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
