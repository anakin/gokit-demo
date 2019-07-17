package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/transport"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

type HelloEndpoints struct {
	HelloEndpoint  endpoint.Endpoint
	HealthEndpoint endpoint.Endpoint
}

type helloReponse struct {
	V   string `json:"s"`
	Err string `json:"err,omitempty"`
}

type helloRequest struct {
	S string `json:"s"`
}
type HealthRequest struct{}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthEndpoint(svc HelloService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{status}, nil
	}
}

func MakeHelloEndpoint(svc HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(helloRequest)
		v, err := svc.Hello(req.S)
		// fmt.Println("v:", v, "err:", err)
		if err != nil {
			return helloReponse{"", err.Error()}, nil
		}
		return helloReponse{v, ""}, nil
	}
}

func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req helloRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

func MakeHttpHandler(ctx context.Context, endpoints HelloEndpoints, logger transport.ErrorHandler) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(logger),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
	}
	r.Methods("POST").Path("/hello").Handler(kithttp.NewServer(
		endpoints.HelloEndpoint,
		decodeHelloRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthEndpoint,
		decodeHealthCheckRequest,
		encodeResponse,
		options...,
	))
	return r
}
