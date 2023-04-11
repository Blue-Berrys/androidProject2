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

type LikeHandler struct {
	*gin.Context
	UserId        uint
	FriendsChatId uint
	actionType    int
}

func LikeActionHandler(c *gin.Context) {
	NewLikeActionHandler(c).Do()
}

func NewLikeActionHandler(c *gin.Context) *LikeHandler {
	return &LikeHandler{Context: c}
}

func (q *LikeHandler) Do() {
	//Get Parameter
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error()) //自定义的ErrorString,将err转化成string
		return
	}

	//use Parameter
	if err := service.LikeAction(q.UserId, q.FriendsChatId, q.actionType); err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk()
}

func (q *LikeHandler) ParseParameter() error {
	wg := sync.WaitGroup{}
	wg.Add(3)

	errChan := make(chan error, 3)
	defer close(errChan)

	go func() {
		defer wg.Done()
		curUserId, _ := q.Get("UserId") //curUserId 是 interface类型的
		UserId, ok := curUserId.(uint)  //转uint
		if !ok {
			errStr := "ParseUserId Failed"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.UserId = UserId
	}()

	go func() {
		defer wg.Done()
		//获取FriendsChatId
		FriendschatIdStr := q.PostForm("friendschat_id")
		FriendsChatId, err := strconv.ParseInt(FriendschatIdStr, 10, 32)
		if err != nil {
			errStr := "ParseFriendsChatId Failed"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.FriendsChatId = uint(FriendsChatId)
	}()

	go func() {
		defer wg.Done()
		//获取actionType
		strActionType := q.PostForm("action_type")
		ActionType, err := strconv.ParseInt(strActionType, 10, 32)
		if err != nil {
			errStr := "ParseActionType Failed"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.actionType = int(ActionType)
	}()

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (q *LikeHandler) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *LikeHandler) SendOk() {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 0,
		StatusMsg:  config.SUCCESS_MSG,
	})
}
