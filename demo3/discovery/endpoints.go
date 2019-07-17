package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/sd"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

func MakeDiscoverEndpoint(ctx context.Context, client consul.Client, logger log.Logger) endpoint.Endpoint {
	serviceName := "hello-service"
	tags := []string{"chope", "b2c"}
	passingOnly := true
	duration := 500 * time.Millisecond

	instancer := consul.NewInstancer(client, logger, serviceName, tags, passingOnly)
	factory := helloFactory(ctx, "POST", "/hello")
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(1, duration, balancer)
	return retry
}
