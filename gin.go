package main

import (
	"dds_core_server/kafka"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	// 初始化kafka
	err := kafka.InitKafkaProducer([]string{"47.98.225.138:9092"})
	if err != nil {
		log.Fatalf("kafka producer init err: %v", err)
		return
	}
	defer kafka.Producer.Close()

	r.GET("/send_message", func(c *gin.Context) {

		err := kafka.Produce(kafka.KafkaMessage{
			Action: "Add",
			IP:     "11.1.1.1",
			Topic:  "test",
		})
		if err != nil {
			c.JSON(500, resp.InternalError(err))
			return
		}

		c.JSON(200, resp.Success())
	})

	r.Run() // 默认监听并在 0.0.0.0:8080 上启动服务
}
