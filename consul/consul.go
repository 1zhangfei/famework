package consul

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"net"
)

func RegisterConsul(addip string, port int64, port1 string, name string) error {
	clin, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%v:%v", addip, port1)})
	if err != nil {
		return err
	}
	ip := GetIp()

	err = clin.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    name,
		Tags:    []string{"GRPC"},
		Port:    int(port),
		Address: ip[0],
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%v:%v", ip[0], port),
			DeregisterCriticalServiceAfter: "30s",
		},
	})

	if err != nil {
		return err
	}
	return nil

}

func FindConsAddress(ip, port, name string) (string, int64, error) {
	clin, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%v:%v", ip, port)})
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

func GetIp() (ip []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, addr := range addrs {
		ipNet, isVailIpNet := addr.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = append(ip, ipNet.IP.String())
			}
		}

	}
	return ip
}
