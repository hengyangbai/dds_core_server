package controller

import (
	"dds_core_server/kafka"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type SendInfoRequest struct {
	IP    string `json:"ip"`
	Topic string `json:"topic"`
}

func SendInfo(c *gin.Context) {
	var req SendInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求参数有问题")
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, resp.InValidRequest(err))
	}
	fmt.Println("receive IP : ", req.IP)
	// 执行send逻辑

	err := kafka.Produce(&kafka.KafkaMessage{
		Action: kafka.ACTION_ADD,
		IP:     req.IP,
		Topic:  "test",
	})
	if err != nil {
		log.Println("Produce error : ", err)
		return
	}

	IPList := strings.Split(req.IP, ".")
	if len(IPList) != 4 {
		c.JSON(http.StatusBadRequest, resp.InValidRequest(errors.New("ip format error")))
	}

	c.JSON(http.StatusOK, resp.Success())
}
