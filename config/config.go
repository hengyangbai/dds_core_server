package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Kafka Kafka `mapstructure:"kafka" json:"kafka"`
	Redis Redis `mapstructure:"redis" json:"redis"`
}

type Kafka struct {
	Addr []string `mapstructure:"addr" json:"addr"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Db       int32  `mapstructure:"db" json:"db"`
	PoolSize int32  `mapstructure:"pool_size" json:"pool_size"`
}

func ConfigInit() *Config {
	config := &Config{}
	viper.SetConfigFile("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	log.Println("config init success", config)
	return config
}
