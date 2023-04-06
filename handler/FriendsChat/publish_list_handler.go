package handler

import (
	"androidProject2/config"
	service "androidProject2/service/FriendsChat"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type PublishListResponses struct {
	CommonResponse
	*service.PublishListResponse
}

type PublishListHandlerStruct struct {
	*gin.Context
	UserId uint
	SeenId uint
}

func PublishListHandler(c *gin.Context) {
	NewPublishListHandler(c).Do()
}

func NewPublishListHandler(c *gin.Context) *PublishListHandlerStruct {
	return &PublishListHandlerStruct{Context: c}
}

func (q *PublishListHandlerStruct) Do() {
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	info, err := service.PublishList(q.UserId, q.SeenId)
	if err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk(info)
}

func (q *PublishListHandlerStruct) ParseParameter() error {
	//获取user_id
	rawUserId, _ := q.Get("UserId")
	log.Println("tokenId: ", rawUserId)
	UserId, ok := rawUserId.(uint)
	q.UserId = UserId
	if !ok {
		return errors.New("ParseUserId Failed") //创建错误
	}
	SeenIdStr := q.PostForm("user_id")
	seenId, err := strconv.ParseInt(SeenIdStr, 10, 64)
	if err != nil {
		return errors.New("传入的user_id不是整型")
	}
	q.SeenId = uint(seenId)
	return nil
}

func (q *PublishListHandlerStruct) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *PublishListHandlerStruct) SendOk(info *service.PublishListResponse) {
	q.JSON(http.StatusOK, PublishListResponses{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		PublishListResponse: info,
	})
}
