package grpc

import (
	"encoding/json"
	"fmt"
	"github.com/1zhangfei/famework/config"
	"github.com/1zhangfei/famework/consul"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Rgc struct {
	Add struct {
		Ip   string
		Port int64
	} `json:"rpc"`
	Name   string `json:"tokenName"`
	Consul string `json:"consul"`
}

func getClient(address string) (*Rgc, error) {
	err := config.ViperInit(address)
	if err != nil {
		return nil, err
	}
	id := viper.GetString("Grpc.DataId")
	Group := viper.GetString("Grpc.Group")
	cnf, err := config.GetConfig(id, Group)
	fmt.Println("11111", cnf)
	if err != nil {
		return nil, err
	}
	var r Rgc
	if err = json.Unmarshal([]byte(cnf), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func Service(address string, register func(s *grpc.Server)) error {
	c, err := getClient(address)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.Add.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = consul.RegisterConsul(c.Add.Ip, c.Add.Port, c.Consul, c.Name)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	//反射机制
	reflection.Register(s)
	//健康检测
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}
