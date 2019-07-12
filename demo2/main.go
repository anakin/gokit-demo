package main

import (
	"github.com/go-kit/kit/log"
	ht "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

func main() {
logger:=log.NewLogfmtLogger(os.Stderr)
var svc HelloService

svc = helloService{}
svc = loggingMiddleware{logger,svc}
helloHandler := ht.NewServer(
	makeHelloEndpoint(svc),
	decodeHelloRequest,
encodeResponse,
	)
http.Handle("/hello",helloHandler)
_=logger.Log("msg","http","addr",":8080")
_=logger.Log("err",http.ListenAndServe(":8080",nil))
}
