package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	ht "github.com/go-kit/kit/transport/http"
	"github.com/prometheus/common/log"
	"net/http"
)

type HelloService interface {
	Hello(string) (string, error)
}

type helloService struct {
}

func (h helloService) Hello(s string) (string, error) {
	if s == "" {
		return "", errors.New("no param")
	}
	return "hello " + s, nil
}

type helloReponse struct {
	V   string `json:"s"`
	Err string `json:"err,omitempty"`
}

type helloRequest struct {
	S string `json:"s"`
}

func makeHelloEndpoint(svc HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(helloRequest)
		v, err := svc.Hello(req.S)
		fmt.Println(v)
		if err != nil {
			return helloReponse{"", err.Error()}, nil
		}
		return helloReponse{v, ""}, nil
	}
}

func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req helloRequest
	fmt.Println("req:", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("request err:", err)
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

func main() {
	svc := helloService{}
	helloHandler := ht.NewServer(
		makeHelloEndpoint(svc),
		decodeHelloRequest,
		encodeResponse,
	)
	http.Handle("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
