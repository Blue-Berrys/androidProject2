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

type CommentActionResponses struct {
	CommonResponse
	*service.CommentActionResponse
}

type CommentActionHandlerStruct struct {
	*gin.Context
	UserId        uint
	FriendsChatId uint
	ActionType    int
	CommentText   string
	CommentId     uint
}

func CommentActionHandler(c *gin.Context) {
	NewCommentActionHandler(c).Do()
}

func NewCommentActionHandler(c *gin.Context) *CommentActionHandlerStruct {
	return &CommentActionHandlerStruct{Context: c}
}

func (q *CommentActionHandlerStruct) Do() {
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	info, err := service.CommentAction(q.UserId, q.FriendsChatId, q.ActionType, q.CommentText, q.CommentId)
	if err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk(info)
}

func (q *CommentActionHandlerStruct) ParseParameter() error {
	wg := sync.WaitGroup{}
	wg.Add(3)

	errChan := make(chan error, 3)
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

	go func() {
		defer wg.Done()
		ActionTypeStr := q.PostForm("action_type")
		ActionType, err := strconv.ParseInt(ActionTypeStr, 10, 64)
		if err != nil {
			errStr := "action_type不是整型"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.ActionType = int(ActionType)
	}()
	wg.Wait()

	if q.ActionType == 1 {
		CommentText := q.PostForm("comment_text")
		q.CommentText = CommentText

		q.CommentId = 0
	} else if q.ActionType == 2 {
		CommentIdStr := q.PostForm("comment_id")
		CommentId, err := strconv.ParseInt(CommentIdStr, 10, 64)
		if err != nil {
			errStr := "comment_id不是整型"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.CommentId = uint(CommentId)

		q.CommentText = ""
	} else {
		errStr := "action_type输入错误"
		log.Println(errStr)
		errChan <- errors.New(errStr)
	}

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (q *CommentActionHandlerStruct) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *CommentActionHandlerStruct) SendOk(info *service.CommentActionResponse) {
	q.JSON(http.StatusOK, CommentActionResponses{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		CommentActionResponse: info,
	})
}
