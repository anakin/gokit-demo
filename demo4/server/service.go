package server

import (
	"context"
)

type UserService interface {
	Get(context.Context, int32) (string, error)
}

type userService struct{}

func New() UserService {
	return userService{}
}

func (userService) Get(_ context.Context, id int32) (string, error) {
	//time.Sleep(time.Second * 2)
	return "userinfo:aaaa", nil
}
