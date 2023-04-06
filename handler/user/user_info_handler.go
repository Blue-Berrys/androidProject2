package handler

import (
	"androidProject2/config"
	service "androidProject2/service/user"
	"errors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
)

type UserInfoResponse struct {
	CommonResponse
	*service.InfoResponse
}

type InfoHandler struct {
	*gin.Context
	userId              uint
	SeenId              uint
	action_type         int
	signature           string
	avatar              *multipart.FileHeader
	background_image    *multipart.FileHeader
	AvatarName          string
	BackgroundImageName string
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

	info, err := service.QueryUserInfo(q.userId, q.SeenId, q.action_type, q.signature, q.avatar, q.background_image, q.AvatarName, q.BackgroundImageName)
	if err != nil {
		q.SendError(err.Error())
		return
	}
	q.SendOk(info)
}

func (q *InfoHandler) ParseParameter() error {
	//获取user_id
	rawUserId, _ := q.Get("UserId")
	log.Println("tokenId: ", rawUserId)
	UserId, ok := rawUserId.(uint)
	if !ok {
		return errors.New("ParseUserId Failed") //创建错误
	}
	q.userId = UserId

	//被看的人id
	SeenIdStr := q.PostForm("user_id")
	SeenId, err := strconv.ParseInt(SeenIdStr, 10, 64)
	if err != nil {
		errors.New("传入的user_id不是整型")
	}
	q.SeenId = uint(SeenId)

	actionTypeStr := q.PostForm("action_type")
	action_type, err := strconv.ParseInt(actionTypeStr, 10, 64)
	q.action_type = int(action_type)
	if err != nil {
		errors.New("传入的action_type不是整型")
	}

	if action_type == 1 {
		var signature string
		var avatar *multipart.FileHeader
		var background_image *multipart.FileHeader
		signature = q.PostForm("signature")
		avatar, err = q.FormFile("avatar")
		if err != nil {
			return err
		}
		background_image, err := q.FormFile("background_image")
		if err != nil {
			return err
		}
		q.background_image = background_image
		q.avatar = avatar
		q.signature = signature
		q.AvatarName = uuid.NewV4().String()
		q.BackgroundImageName = uuid.NewV4().String()
		avatardst := path.Join("~/imageorepo", q.AvatarName+".jpg")
		background_imagedst := path.Join("~/imageorepo", q.BackgroundImageName+".jpg")
		q.SaveUploadedFile(q.avatar, avatardst)
		q.SaveUploadedFile(q.background_image, background_imagedst)
	} else {
		q.background_image = nil
		q.avatar = nil
		q.signature = ""
		q.BackgroundImageName = ""
		q.AvatarName = ""
	}
	log.Println("token_id:", UserId, "SeenId", SeenId)
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
