package controller

import (
	"context"
	"dds_core_server/kafka"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type Repo interface {
	AddValue(ctx context.Context, info *BaseInfo) error
	DelValue(ctx context.Context, info *BaseInfo) error
	DeleteAll(ctx context.Context) error
	GetInfoList(ctx context.Context, pageIndex int32, pageSize int32) ([]BaseInfo, int64, error)
}

type Controller struct {
	repo Repo
}

func NewController(repo Repo) *Controller {
	return &Controller{repo: repo}
}

type SendInfoRequest struct {
	IP    string `json:"ip"`
	Topic string `json:"topic"`
	Port  int32  `json:"port"`
}

func (ct *Controller) SendInfo(c *gin.Context) {
	var req SendInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求参数有问题")
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, resp.InValidRequest(err))
	}
	fmt.Println("receive IP : ", req.IP)
	IPList := strings.Split(req.IP, ".")
	if len(IPList) != 4 {
		c.JSON(http.StatusBadRequest, resp.InValidRequest(errors.New("ip format error")))
		return
	}

	// 执行send逻辑
	err := kafka.Produce(&kafka.KafkaMessage{
		Action: kafka.ACTION_ADD,
		IP:     req.IP,
		Topic:  "test",
	})
	if err != nil {
		log.Println("Produce error : ", err)
		c.JSON(http.StatusInternalServerError, resp.InternalError(err))
	}

	// 执行存储逻辑
	err = ct.repo.AddValue(c, &BaseInfo{
		IP:    req.IP,
		Topic: req.Topic,
		Port:  req.Port,
	})
	if err != nil {
		log.Println("AddValue error : ", err)
		c.JSON(http.StatusInternalServerError, resp.InternalError(err))
		return
	}

	c.JSON(http.StatusOK, resp.Success())
	return
}

type GetInfoRequest struct {
	PageSize  int32 `json:"page_size"`
	PageIndex int32 `json:"page_index"`
}

type GetInfoResponse struct {
	Data    *InfoData `json:"data"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
}
type InfoData struct {
	BaseInfos []BaseInfo `json:"base_infos"`
	PageIndex int32      `json:"page_index"`
	PageSize  int32      `json:"page_size"`
	Total     int32      `json:"total"`
}

type BaseInfo struct {
	IP         string `json:"ip"`
	Topic      string `json:"topic"`
	Port       int32  `json:"port"`
	CreateTime int64  `json:"create_time"`
}

func (ct *Controller) GetInfoList(c *gin.Context) {
	var req GetInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Println("请求参数有问题")
		log.Println(err)
		c.JSON(http.StatusBadRequest, resp.InValidRequest(err))
	}
	fmt.Println("PageSize : ", req.PageSize)
	fmt.Println("PageIndex : ", req.PageIndex)
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var response GetInfoResponse
	//response.Data = make([]BaseInfo, 0)
	response.Code = 0
	response.Message = "success"

	// redis 查询列表数据
	list, count, err := ct.repo.GetInfoList(c, req.PageIndex, req.PageSize)
	if err != nil {
		log.Println("GetInfoList error : ", err)
		c.JSON(http.StatusInternalServerError, resp.InternalError(err))
	}
	response.Data = &InfoData{
		BaseInfos: list,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
		Total:     int32(count),
	}

	c.JSON(http.StatusOK, response)
}
