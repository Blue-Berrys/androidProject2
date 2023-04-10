package handler

import (
	"androidProject2/config"
	service "androidProject2/service/Comment"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type CommentListResponses struct {
	CommonResponse
	*service.CommentListResponse
}

type CommentListHandlerStruct struct {
	*gin.Context
	UserId        uint
	FriendsChatId uint
}

func CommentListHandler(c *gin.Context) {
	NewCommentListHandler(c).Do()
}

func NewCommentListHandler(c *gin.Context) *CommentListHandlerStruct {
	return &CommentListHandlerStruct{Context: c}
}

func (q *CommentListHandlerStruct) Do() {
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	info, err := service.CommentList(q.UserId, q.FriendsChatId)
	if err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk(info)
}

func (q *CommentListHandlerStruct) ParseParameter() error {
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
		if !ok {
			errStr := "ParseUserId Failed"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.UserId = UserId
	}()

	go func() {
		defer wg.Done()
		FriendsChatIdStr := q.PostForm("friendschat_id")
		FriendsChatId, err := strconv.ParseInt(FriendsChatIdStr, 10, 64)
		if err != nil {
			errStr := "friendschat_id不是整型"
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

func (q *CommentListHandlerStruct) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *CommentListHandlerStruct) SendOk(info *service.CommentListResponse) {
	q.JSON(http.StatusOK, CommentListResponses{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		CommentListResponse: info,
	})
}
