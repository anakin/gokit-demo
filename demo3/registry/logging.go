package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next HelloService) HelloService {
		return logMW{logger, next}
	}
}

//logMW log middleware
type logMW struct {
	logger log.Logger
	HelloService
}

func (mw logMW) Hello(s string) (output string, err error) {
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
