package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next HelloService) HelloService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	HelloService
}

func (mw logmw) Hello(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "hello",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.HelloService.Hello(s)
	return
}
