package handler

import (
	"androidProject2/config"
	service "androidProject2/service/Like"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type LikeListResponse struct {
	CommonResponse
	*service.LikeListsResponse
}

type LikeListHandlerStruct struct {
	*gin.Context
	UserId        uint
	FriendsChatId uint
}

func LikeListHandler(c *gin.Context) {
	NewLikeListHandler(c).Do()
}

func NewLikeListHandler(c *gin.Context) *LikeListHandlerStruct {
	return &LikeListHandlerStruct{Context: c}
}

func (q *LikeListHandlerStruct) Do() {
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	info, err := service.LikeList(q.UserId, q.FriendsChatId)
	if err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk(info)
}

func (q *LikeListHandlerStruct) ParseParameter() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 2)
	defer close(errChan)

	go func() {
		defer wg.Done()
		//获取user_id
		rawUserId, _ := q.Get("UserId")
		log.Println("tokenId: ", rawUserId)
		UserId, ok := rawUserId.(uint)
		q.UserId = UserId
		if !ok {
			errStr := "ParseUserId Failed"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	go func() {
		defer wg.Done()
		FriendsChatIdStr := q.PostForm("friendschat_id")
		FriendsChatId, err := strconv.ParseInt(FriendsChatIdStr, 10, 64)
		if err != nil {
			errStr := "传入的friendschat_id不是整型"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.FriendsChatId = uint(FriendsChatId)
	}()

	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (q *LikeListHandlerStruct) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *LikeListHandlerStruct) SendOk(info *service.LikeListsResponse) {
	q.JSON(http.StatusOK, LikeListResponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		LikeListsResponse: info,
	})
}
