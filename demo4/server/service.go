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
	return "userinfo:aaaa", nil
}
