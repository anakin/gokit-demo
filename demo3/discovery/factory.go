package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
)

type helloRequest struct {
	S string `json:"s"`
}

type helloResponse struct {
	V   string `json:"s"`
	Err error  `json:"err"`
}

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
	fmt.Println("dis req:", request)
	cRequest := request.(CalRequest)
	// jsonStr := []byte(`{"s": "aa"}`)
	// var err error
	// req, err = http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	// if err != nil {
	// 	return err
	// }
	var buf bytes.Buffer
	reqq := helloRequest{S: cRequest.Str}
	if err := json.NewEncoder(&buf).Encode(reqq); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeHelloResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response helloResponse
	var s map[string]interface{}
	if respCode := resp.StatusCode; respCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}
		return nil, errors.New(s["err"].(string) + "\n")
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
