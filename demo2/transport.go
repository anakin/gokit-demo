package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

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

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}
