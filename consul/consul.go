package consul

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

func RegisterConsul(address string, port int64, name string) error {
	clin, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	err = clin.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    name,
		Tags:    []string{"GRPC"},
		Port:    int(port),
		Address: address,
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%v:%v", address, port),
			DeregisterCriticalServiceAfter: "30s",
		},
	})

	if err != nil {
		return err
	}
	return nil

}

func FindConsAddress(name string) (string, int64, error) {
	clin, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", 0, err
	}

	byname, data, err := clin.Agent().AgentHealthServiceByName(name)
	if err != nil {
		return "", 0, err
	}
	if len(data) == 0 {
		return "", 0, errors.New("没有健康的服务")
	}
	if byname == "passing`" {
		return "", 0, errors.New("没有健康的服务")
	}

	return data[0].Service.Address, int64(data[0].Service.Port), nil

}
