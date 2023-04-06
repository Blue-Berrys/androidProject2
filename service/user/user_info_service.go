package service

import (
	"androidProject2/config"
	"androidProject2/middleware/minio"
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
)

type InfoResponse struct {
	MyUser *util.User `json:"user"`
}

type QueryUserInfoFlow struct {
	UserId              uint
	SeenId              uint
	action_type         int
	signature           string
	avatar              *multipart.FileHeader
	background_image    *multipart.FileHeader
	AvatarName          string
	BackgroundImageName string
	modify              bool
	data                *util.User
	*InfoResponse
}

func QueryUserInfo(UserId uint, SeenId uint, action_type int, signature string, avatar *multipart.FileHeader, background_image *multipart.FileHeader, AvatarName string, BackgroundImageName string) (*InfoResponse, error) {
	return NewQueryUserInfo(UserId, SeenId, action_type, signature, avatar, background_image, AvatarName, BackgroundImageName).Do()
}

func NewQueryUserInfo(Userid uint, SeenId uint, action_type int, signature string, avatar *multipart.FileHeader, background_image *multipart.FileHeader, AvatarName string, BackgroundImageName string) *QueryUserInfoFlow {
	return &QueryUserInfoFlow{UserId: Userid, SeenId: SeenId, action_type: action_type, signature: signature, avatar: avatar, background_image: background_image, AvatarName: AvatarName, BackgroundImageName: BackgroundImageName}
}

func (q *QueryUserInfoFlow) Do() (*InfoResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.InfoResponse, nil
}

func (q *QueryUserInfoFlow) checkNum() error {
	//查action_type是否只有1或2
	if q.action_type != 1 && q.action_type != 2 {
		return errors.New("action_type不是1和2，输入错误")
	}
	//查询用户id是否在数据库中
	var UserDao = model.NewUserDao()
	fmt.Println(q.UserId, q.SeenId)
	if q.UserId == 0 || !UserDao.QueryUserExistByUserId(q.UserId) {
		return errors.New("UserId用户不存在")
	}
	if q.SeenId == 0 || !UserDao.QueryUserExistByUserId(q.SeenId) {
		return errors.New("被看的用户不存在")
	}
	q.modify = false
	log.Println("action_type:", q.UserId, q.SeenId, q.action_type)
	//自己看自己并且要修改
	if q.UserId == q.SeenId && q.action_type == 1 {
		q.modify = true
	}
	return nil
}

func (q *QueryUserInfoFlow) prepareData() error {
	//根据id查询用户信息
	var UserDao = model.NewUserDao()
	//保存到数据库
	var dbUser = model2.User{}
	log.Println("modify:", q.modify)
	//修改
	if q.modify {
		if err := minio.ImageToMinio(q.avatar, q.AvatarName); err != nil {
			return err
		}
		if err := minio.ImageToMinio(q.background_image, q.BackgroundImageName); err != nil {
			return err
		}
		log.Println("9999")
		if err := UserDao.UpdateUserInfo(&dbUser, q.SeenId, q.signature, q.AvatarName, q.BackgroundImageName); err != nil {
			return err
		}
	}

	if err := UserDao.QueryUserInfoById(q.SeenId, &dbUser); err != nil {
		return err
	}

	q.data = &util.User{
		Id:              dbUser.ID,
		Name:            dbUser.UserName,
		Signature:       dbUser.Signature,
		WorkCount:       dbUser.WorkCount,
		BackGroundImage: config.Miniourl + dbUser.BackgroundImage,
		Avatar:          config.Miniourl + dbUser.Avatar,
	}

	return nil
}

func (q *QueryUserInfoFlow) packData() error {
	q.InfoResponse = &InfoResponse{MyUser: q.data}
	return nil
}
