package reg

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

type ConsulRegister struct {
	ConsulAddress                  string // consul address
	ServiceName                    string // service name
	ServiceIP                      string
	Tags                           []string // consul tags
	ServicePort                    int      //service port
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func NewConsulRegister(consulAddress, serviceName, serviceIP string, servicePort int, tags []string) *ConsulRegister {
	return &ConsulRegister{
		ConsulAddress:                  consulAddress,
		ServiceName:                    serviceName,
		ServiceIP:                      serviceIP,
		Tags:                           tags,
		ServicePort:                    servicePort,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

// https://github.com/ru-rocker/gokit-playground/blob/master/lorem-consul/register.go
// https://github.com/hatlonely/hellogolang/blob/master/sample/addservice/internal/grpcsr/consul_register.go
func (r *ConsulRegister) NewConsulGRPCRegister() (*consulsd.Registrar, error) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	consulConfig := api.DefaultConfig()
	consulConfig.Address = r.ConsulAddress
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	client := consulsd.NewClient(consulClient)

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", r.ServiceName, r.ServiceIP, r.ServicePort),
		Name:    fmt.Sprintf("grpc.health.v1.%v", r.ServiceName),
		Tags:    r.Tags,
		Port:    r.ServicePort,
		Address: r.ServiceIP,
		Check: &api.AgentServiceCheck{
			// 健康检查间隔
			Interval: r.Interval.String(),
			//grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			GRPC: fmt.Sprintf("%v:%v/%v", r.ServiceIP, r.ServicePort, r.ServiceName),
			// 注销时间，相当于过期时间
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),
		},
	}
	return consulsd.NewRegistrar(client, reg, logger), nil
}

func (r *ConsulRegister) NewConsulHttpRegister() (*consulsd.Registrar, error) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	consulCfg := api.DefaultConfig()
	consulCfg.Address = r.ConsulAddress
	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	client := consulsd.NewClient(consulClient)

	// 设置Consul对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP:     "http://" + r.ServiceIP + ":" + strconv.Itoa(r.ServicePort) + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	//设置微服务想Consul的注册信息
	regger := &api.AgentServiceRegistration{
		ID:      r.ServiceName + time.Now().String(),
		Name:    r.ServiceName,
		Address: r.ServiceIP,
		Port:    r.ServicePort,
		Tags:    r.Tags,
		Check:   &check,
	}

	return consulsd.NewRegistrar(client, regger, logger), nil
}
