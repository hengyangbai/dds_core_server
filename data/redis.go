package data

import (
	"context"
	"dds_core_server/config"
	"dds_core_server/controller"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type InfoRepo struct {
	redisClient *redis.Client
}

const Key = "strategy"

func NewRepo(conf config.Redis) (*InfoRepo, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		DB:       int(conf.Db),       // 数据库
		PoolSize: int(conf.PoolSize), // 连接池大小
	})
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis ping err:", err)
		return nil, err
	}
	log.Println(pong)
	return &InfoRepo{redisClient: redisClient}, nil
}
func (rp *InfoRepo) AddValue(ctx context.Context, info *controller.BaseInfo) error {
	// json串 zset
	jsonByte, err := json.Marshal(info)
	if err != nil {
		log.Println("json marshal err:", err)
		return err
	}

	return rp.redisClient.ZAdd(ctx, Key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: string(jsonByte),
	}).Err()

}
func (rp *InfoRepo) DelValue(ctx context.Context, info *controller.BaseInfo) error {
	jsonByte, err := json.Marshal(info)
	if err != nil {
		log.Println("json marshal err:", err)
		return err
	}

	// 使用 ZREM 删除成员
	return rp.redisClient.ZRem(ctx, Key, string(jsonByte)).Err()
}
func (rp *InfoRepo) DeleteAll(ctx context.Context) error {
	return rp.redisClient.Del(ctx, Key).Err()
}
func (rp *InfoRepo) GetInfoList(ctx context.Context, pageIndex int32, pageSize int32) ([]controller.BaseInfo, int64, error) {
	totalCount, err := rp.redisClient.ZCard(ctx, Key).Result()
	if err != nil {
		log.Println("redis zcard err:", err)
		return nil, 0, err
	}

	start := (pageIndex - 1) * pageSize
	stop := start + pageSize - 1

	// 使用 ZRANGE 获取范围内的成员
	result, err := rp.redisClient.ZRevRangeWithScores(ctx, Key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Println("redis zrange err:", err)
		return nil, 0, err
	}

	// 将结果反序列化回 BaseInfo 对象数组
	var infoList []controller.BaseInfo
	for _, z := range result {
		var info controller.BaseInfo
		err := json.Unmarshal([]byte(z.Member.(string)), &info)
		if err != nil {
			log.Println("json unmarshal err:", err)
			continue
		}

		info.CreateTime = int64(z.Score)
		infoList = append(infoList, info)
	}

	return infoList, totalCount, nil
}
