package middleware

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
)

func Hystrix(commandName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			if err := hystrix.Do(commandName, func() (err error) {
				resp, err = next(ctx, request)
				return err
			}, nil); err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
}
