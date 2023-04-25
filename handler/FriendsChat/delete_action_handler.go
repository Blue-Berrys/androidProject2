package handler

import (
	"androidProject2/config"
	service "androidProject2/service/FriendsChat"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type DeleteFriendChatsResponses struct {
	CommonResponse
}

type InfoDeleteFriendChats struct {
	*gin.Context
	userId        uint
	FriendsChatId uint
}

func DeleteFriendChatsHandler(c *gin.Context) {
	NewInfoDeleteFriendsChats(c).Do()
}

func NewInfoDeleteFriendsChats(c *gin.Context) *InfoDeleteFriendChats {
	return &InfoDeleteFriendChats{Context: c}
}

func (q *InfoDeleteFriendChats) Do() {
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	if err := service.DeleteAction(q.userId, q.FriendsChatId); err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk()
}

func (q *InfoDeleteFriendChats) ParseParameter() error {
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
			errChan <- errors.New(errStr) //创建错误
		}
		q.userId = UserId
	}()

	go func() {
		defer wg.Done()
		//获取friendschat_id
		friendsChatIdStr := q.PostForm("friendschat_id")
		log.Println("friendsChatIdStr: ", friendsChatIdStr)
		friendsChatId, err := strconv.ParseInt(friendsChatIdStr, 10, 64)
		if err != nil {
			errStr := "friendschat_id不是整型"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.FriendsChatId = uint(friendsChatId)
	}()

	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (q *InfoDeleteFriendChats) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *InfoDeleteFriendChats) SendOk() {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 0,
		StatusMsg:  config.SUCCESS_MSG,
	})
}
