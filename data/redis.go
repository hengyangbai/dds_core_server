package data

import (
	"context"
	"dds_core_server/config"
	"dds_core_server/controller"
	"github.com/redis/go-redis/v9"
	"log"
)

var redisClient *redis.Client

const Key = "strategy"

func RedisInit(conf config.Redis) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		DB:       int(conf.Db),       // 数据库
		PoolSize: int(conf.PoolSize), // 连接池大小
	})
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis ping err:", err)
		return err
	}
	log.Println(pong)
	return nil
}

func AddValue(info controller.BaseInfo) {
	// json串

}

func DelValue(value string) {

}

func DeleteAll() {

}
