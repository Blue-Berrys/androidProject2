package handler

import (
	"androidProject2/config"
	service "androidProject2/service/user"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserInfoResponse struct {
	CommonResponse
	*service.InfoResponse
}

type InfoHandler struct {
	*gin.Context
	userId uint
}

func UserInfoHandler(c *gin.Context) {
	NewUserInfoHandler(c).Do()
}

func NewUserInfoHandler(c *gin.Context) *InfoHandler {
	return &InfoHandler{Context: c}
}

func (q *InfoHandler) Do() {
	//获得参数
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	info, err := service.QueryUserInfo(q.userId)
	if err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk(info)
}

func (q *InfoHandler) ParseParameter() error {
	//获取user_id
	rawUserId, _ := q.Get("UserId")
	UserId, ok := rawUserId.(uint)
	if !ok {
		return errors.New("ParseUserId Failed") //创建错误
	}
	q.userId = UserId
	return nil
}

func (q *InfoHandler) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *InfoHandler) SendOk(info *service.InfoResponse) {
	q.JSON(http.StatusOK, UserInfoResponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		InfoResponse: info,
	})
}
