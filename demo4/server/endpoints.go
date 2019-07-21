package server

import (
	"context"
	pb "demo4/proto/user"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetEndpoint endpoint.Endpoint
}

type userRequest struct {
	Id int32
}

type userResponse struct {
	Message string
	Err     error
}

func MakeUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*pb.UserRequest)
		v, err := svc.Get(ctx, r.Userid)
		if err != nil {

			return userResponse{Message: "", Err: err}, nil
		}
		return userResponse{Message: v, Err: nil}, nil
	}
}

func (e Endpoints) Get(ctx context.Context, id int32) (string, error) {
	req := userRequest{Id: id}
	fmt.Println("request here:", req)
	res, err := e.GetEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	return res.(userResponse).Message, nil
}
