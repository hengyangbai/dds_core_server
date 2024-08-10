package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
)

const TOPIC = "web_log"

var Producer sarama.SyncProducer

func InitKafkaProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	var err error
	// 连接kafka
	Producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("producer closed, err: %v", err)
	}

	// 关闭 producer 的操作应该在你的程序退出时发生，比如可以在 main 函数的 defer 中调用
	// defer Producer.Close()
	return nil
}
