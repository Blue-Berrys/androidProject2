package handler

import (
	"androidProject2/config"
	service "androidProject2/service/user"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type LevelActionFlow struct {
	*gin.Context
	UserId       uint
	ModifyUserId uint
}

func LevelActionHandler(c *gin.Context) {
	NewLevelActionHandler(c).Do()
}

func NewLevelActionHandler(c *gin.Context) *LevelActionFlow {
	return &LevelActionFlow{Context: c}
}

func (q *LevelActionFlow) Do() {
	//获得参数
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	if err := service.LevelActionService(q.UserId, q.ModifyUserId); err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk()
}

func (q *LevelActionFlow) ParseParameter() error {
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
		q.UserId = UserId
	}()

	go func() {
		defer wg.Done()
		//传入的user_id
		ModifyUserIdStr := q.PostForm("user_id")
		log.Println("ModifyUserIdStr:", ModifyUserIdStr)
		ModifyUserId, err := strconv.ParseInt(ModifyUserIdStr, 10, 64)
		if err != nil {
			errStr := "传入的user_id不是整型"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.ModifyUserId = uint(ModifyUserId)
	}()
	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (q *LevelActionFlow) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *LevelActionFlow) SendOk() {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 0,
		StatusMsg:  config.SUCCESS_MSG,
	})
}
