package server

import (
	"context"
	pb "demo4/proto/user"
	"errors"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get grpctransport.Handler
}

func MakeGRPCServer(endpoints Endpoints) pb.UserServer {
	return &grpcServer{
		get: grpctransport.NewServer(
			endpoints.GetEndpoint,
			DecodeGRPCGetRequest,
			EncodeGRPCGetRssponse,
		),
	}
}

func (s *grpcServer) Get(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	_, resp, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserResponse), nil
}

func DecodeGRPCGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserRequest)
	fmt.Println("server request:", req)
	return &pb.UserRequest{
		Userid:   req.Userid,
		Username: "wwww",
	}, nil
}

func EncodeGRPCGetRssponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(userResponse)
	if resp.Err != nil {
		return &pb.UserResponse{Message: resp.Message, Err: resp.Err.Error()}, nil
	} else {
		return &pb.UserResponse{Message: resp.Message, Err: ""}, nil
	}
}

func EncodeGRPCGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(userRequest)
	return &pb.UserRequest{Userid: req.Id, Username: "aa"}, nil
}

func DecodeGRPCGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserResponse)
	return userResponse{Message: resp.Message, Err: errors.New(resp.Err)}, nil
}
