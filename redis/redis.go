package redis

import (
	"encoding/json"
	"fmt"
	"github.com/1zhangfei/famework/config"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

func WithRedisCli(address string, hand func(cli *redis.Client) (string, error)) (string, error) {
	err := config.ViperInit(address)
	if err != nil {
		return "", err
	}

	id := viper.GetString("Database.DataId")
	Group := viper.GetString("Database.Group")
	fmt.Println("122123", id, Group)
	type RedisConf struct {
		Host string
		Port string
	}
	var val struct {
		R RedisConf `json:"Redis"`
	}

	res, err2 := config.GetConfig(id, Group)

	if err2 != nil {
		return "", err
	}
	if err = json.Unmarshal([]byte(res), &val); err != nil {
		return "", err
	}

	cli := redis.NewClient(&redis.Options{
		Addr:     val.R.Host + ":" + val.R.Port,
		Password: "",
		DB:       1,
	})

	add, err := hand(cli)
	if err != nil {
		return "", err
	}

	defer func(cli *redis.Client) {
		if err = cli.Close(); err != nil {
			log.Println("***********redis关闭失败=========")
		}
	}(cli)
	return add, nil
}
