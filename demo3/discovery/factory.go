package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
)

func helloFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, err error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		return kithttp.NewClient(method, tgt, encodeHelloRequest, decodeHelloResponse).Endpoint(), nil, nil
	}
}

func encodeHelloRequest(_ context.Context, req *http.Request, request interface{}) error {
	helloRequest := request.(HelloRequest)
	return nil
}

func decodeHelloResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response HelloResponse
	var s map[string]interface{}
	if respCode := resp.StatusCode; respCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}
		return nil, errors.New(s["error"].(string) + "\n")
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
