package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

//CalRequest request struct
type CalRequest struct {
	Str string `json:"str"`
}

func decodeDiscoverRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeDiscoverResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHttpHandler(endpoint endpoint.Endpoint) http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/calculate").Handler(kithttp.NewServer(
		endpoint,
		decodeDiscoverRequest,
		encodeDiscoverResponse,
	))

	return r
}
