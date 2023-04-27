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

type LevelListResponses struct {
	CommonResponse
	*service.UserLevelResponse
}

type LevelListFlow struct {
	*gin.Context
	UserId uint
	Level  int
}

func LevelListHandler(c *gin.Context) {
	NewLevelListHandler(c).Do()
}

func NewLevelListHandler(c *gin.Context) *LevelListFlow {
	return &LevelListFlow{Context: c}
}

func (q *LevelListFlow) Do() {
	//获得参数
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}

	info, err := service.QueryLevelList(q.UserId, q.Level)
	if err != nil {
		q.SendError(err.Error())
		return
	}

	q.SendOk(info)
}

func (q *LevelListFlow) ParseParameter() error {
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
		//传入的level
		LevelStr := q.PostForm("level")
		Level, err := strconv.ParseInt(LevelStr, 10, 64)
		if err != nil {
			errStr := "传入的level不是整型"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
		q.Level = int(Level)
	}()
	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (q *LevelListFlow) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *LevelListFlow) SendOk(info *service.UserLevelResponse) {
	q.JSON(http.StatusOK, LevelListResponses{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		UserLevelResponse: info,
	})
}
