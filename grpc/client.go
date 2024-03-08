package grpc

import (
	"fmt"
	"github.com/1zhangfei/famework/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(address string) (*grpc.ClientConn, error) {
	c, err2 := getClient(address)
	if err2 != nil {
		return nil, err2
	}
	consAddress, port, err := consul.FindConsAddress(c.Name)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("%v:%v", consAddress, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

}
