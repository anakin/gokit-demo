package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingMiddleware struct {
	logger log.Logger
	next   HelloService
}

func (mw loggingMiddleware) Hello(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "hello",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.next.Hello(s)
	return
}
