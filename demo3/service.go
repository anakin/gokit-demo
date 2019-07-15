package main

import "errors"

type HelloService interface {
	Hello(string) (string, error)
	HealthCheck() bool
}

type helloService struct {
}

func (helloService) Hello(s string) (string, error) {
	if s == "" {
		return "", errors.New("param is empty")
	}
	return "hello " + s, nil
}
func (helloService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(HelloService) HelloService
