package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/1zhangfei/famework/config"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//var (
//	Mdb *gorm.DB
//)
//
//func MysqlInit(address string) error {
//	err := config.ViperInit(address)
//	if err != nil {
//		return err
//	}
//	ip := viper.GetString("Grpc.Ip")
//	Group := viper.GetString("Grpc.Group")
//
//	cnf, err := config.GetConfig(ip, Group)
//	if err != nil {
//		return err
//	}
//	var val vals
//	if err = json.Unmarshal([]byte(cnf), &val); err != nil {
//		return err
//	}
//
//	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
//	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
//		val.M.Username,
//		val.M.Password,
//		val.M.Host,
//		val.M.Port,
//		val.M.Database,
//	)
//	Mdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//
//	return nil
//}

func WithMysqlClient(address string, hand func(cli *gorm.DB) error) error {
	err := config.ViperInit(address)
	if err != nil {
		return err
	}
	ip := viper.GetString("Database.DataIp")
	Group := viper.GetString("Database.Group")
	type MysqlConf struct {
		Username string
		Password string
		Host     string
		Port     string
		Database string
	}

	var val struct {
		M MysqlConf `json:"Mysql"`
	}

	cnf, err := config.GetConfig(ip, Group)
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(cnf), &val); err != nil {
		return err
	}

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		val.M.Username,
		val.M.Password,
		val.M.Host,
		val.M.Port,
		val.M.Database,
	)
	db, errs := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Db, _ := db.DB()
	defer func(Db *sql.DB) {
		err = Db.Close()
		if err != nil {
			fmt.Println("*********mysql关闭链接失败，*************")
		}
	}(Db)

	if errs != nil {
		return err
	}
	if err = hand(db); err != nil {
		return err
	}
	return nil
}
