package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var client consul.Client
	cfg := api.DefaultConfig()
	cfg.Address = "172.16.21.244:8500"
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		_ = logger.Log("err", err)
		os.Exit(1)
	}
	client = consul.NewClient(consulClient)

	ctx := context.Background()

	dEnpoint := MakeDiscoverEndpoint(ctx, client, logger)
	r := MakeHttpHandler(dEnpoint)
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		_ = logger.Log("transport", "HTTP", "addr", "9001")
		errChan <- http.ListenAndServe(":9001", r)
	}()
	_ = logger.Log("exit", <-errChan)
}
