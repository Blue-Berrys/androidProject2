package handler

import (
	"androidProject2/config"
	service "androidProject2/service/FriendsChat"
	"errors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
)

type PublishActionResponses struct {
	CommonResponse
	*service.PublishActionResponse
}

type PublishActionHandlerStruct struct {
	*gin.Context
	UserId     uint
	Content    string
	ActionType int
	Images     []*multipart.FileHeader
	ImageName  []string
}

func PublishActionHandler(c *gin.Context) {
	NewPublishActionHandler(c).Do()
}

func NewPublishActionHandler(c *gin.Context) *PublishActionHandlerStruct {
	return &PublishActionHandlerStruct{Context: c}
}

func (q *PublishActionHandlerStruct) Do() {
	if err := q.ParseParameter(); err != nil {
		q.SendError(err.Error())
		return
	}
	info, err := service.PublishAction(q.UserId, q.Content, q.ActionType, q.ImageName, q.Images)
	if err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk(info)
}

func (q *PublishActionHandlerStruct) ParseParameter() error {
	//获取user_id
	rawUserId, _ := q.Get("UserId")
	log.Println("tokenId: ", rawUserId)
	UserId, ok := rawUserId.(uint)
	q.UserId = UserId
	if !ok {
		return errors.New("ParseUserId Failed") //创建错误
	}

	content := q.PostForm("content")
	q.Content = content

	actionTypeStr := q.PostForm("action_type")
	action_type, err := strconv.ParseInt(actionTypeStr, 10, 64)
	q.ActionType = int(action_type)
	if err != nil {
		return errors.New("传入的action_type不是整型")
	}

	//有图片
	if q.ActionType == 1 {
		form, _ := q.MultipartForm()
		images := form.File["image"]
		for i, image := range images {
			log.Println(i)
			q.ImageName = append(q.ImageName, uuid.NewV4().String())
			imagedst := path.Join("~/imageorepo", q.ImageName[i]+".jpg")
			q.SaveUploadedFile(image, imagedst)
		}
		q.Images = images
	} else { //没图片
		q.Images = nil
	}
	return nil
}

func (q *PublishActionHandlerStruct) SendError(msg string) {
	q.JSON(http.StatusOK, CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (q *PublishActionHandlerStruct) SendOk(info *service.PublishActionResponse) {
	q.JSON(http.StatusOK, PublishActionResponses{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		PublishActionResponse: info,
	})
}
