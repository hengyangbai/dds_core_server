package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type KafkaMessage struct {
	Action string `json:"action"`
	IP     string `json:"ip"`
	Topic  string `json:"topic"`
}

const (
	ACTION_ADD    = "Add"
	ACTION_DELETE = "Delete"
)

func Produce(msg KafkaMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("json marshal msg err:%v", err)
		return err
	}
	// 构造一个消息
	msgDes := &sarama.ProducerMessage{
		Topic: TOPIC,
		Value: sarama.StringEncoder(msgBytes),
	}

	// 发送消息
	pid, offset, err := Producer.SendMessage(msgDes)
	if err != nil {
		return fmt.Errorf("send msg failed, err: %v", err)
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
	return nil
}
