package app

import (
	"github.com/1zhangfei/famework/config"
	"github.com/spf13/viper"
)

func Init(address string) error {

	if err := config.ViperInit(address); err != nil {
		return err
	}
	ip := viper.GetString("Nacos.Ip")
	Port := viper.GetInt("Nacos.Port")

	err := config.GetClient(ip, int64(Port))
	if err != nil {
		return err
	}

	return nil
}
