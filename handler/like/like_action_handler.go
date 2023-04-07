package handler

import (
	"androidProject2/config"
	service "androidProject2/service/Like"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	curUserId, _ := q.Get("UserId") //curUserId 是 interface类型的
	UserId, ok := curUserId.(uint)  //转uint
	if !ok {
		return errors.New("ParseUserId Failed") //创建错误
	}
	//获取FriendsChatId
	FriendschatIdStr := q.PostForm("friendschat_id")
	FriendsChatId, err := strconv.ParseInt(FriendschatIdStr, 10, 32)
	if err != nil {
		return errors.New("ParseFriendsChatId Failed")
	}
	//获取actionType
	strActionType := q.PostForm("action_type")
	ActionType, err := strconv.ParseInt(strActionType, 10, 32)
	if err != nil {
		return errors.New("ParseActionType Failed")
	}
	q.actionType = int(ActionType)
	q.FriendsChatId = uint(FriendsChatId)
	q.UserId = UserId
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
