package main

import "errors"

type HelloService interface {
	Hello(string) (string, error)
}

type helloService struct {
}

func (helloService) Hello(s string) (string, error) {
	if s == "" {
		return "", errors.New("param is empty")
	}
	return "hello " + s, nil
}
