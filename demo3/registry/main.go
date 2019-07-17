package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/transport"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"demo3/registry/reg"
	"os/signal"

	"github.com/go-kit/kit/log"
)

func main() {
	var (
		consulHost  = flag.String("consul.host", "172.16.21.244:8500", "consul ip address")
		serviceHost = flag.String("service.host", "172.16.30.26", "service ip address")
		servicePort = flag.Int("service.port", 8080, "service port")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logHandler := transport.NewLogErrorHandler(logger)
	var svc HelloService
	svc = helloService{}
	svc = loggingMiddleware(logger)(svc)
	helloHandler := MakeHelloEndpoint(svc)
	healthHandler := MakeHealthEndpoint(svc)
	epts := HelloEndpoints{
		HelloEndpoint:  helloHandler,
		HealthEndpoint: healthHandler,
	}

	handler := MakeHttpHandler(context.Background(), epts, logHandler)
	conf := reg.NewConsulRegister(*consulHost, "hello-service", *serviceHost, *servicePort, []string{"chope", "b2c"})

	register, err := conf.NewConsulHttpRegister()
	if err != nil {
		fmt.Println(err)
	}
	errChan := make(chan error)

	go func() {
		register.Register()
		r := handler
		_ = logger.Log("msg", "http", "addr", ":8080")
		host := ":" + strconv.Itoa(*servicePort)
		errChan <- http.ListenAndServe(host, r)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	err = <-errChan
	register.Deregister()
	fmt.Println(err)
}
