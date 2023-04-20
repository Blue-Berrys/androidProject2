package service

import (
	"androidProject2/cache/minio"
	"androidProject2/config"
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"log"
	"mime/multipart"
	"sync"
)

type InfoResponse struct {
	MyUser *util.User `json:"user"`
}

type QueryUserInfoFlow struct {
	UserId              uint
	SeenId              uint
	action_type         int
	nickname            string
	avatar              *multipart.FileHeader
	background_image    *multipart.FileHeader
	AvatarName          string
	BackgroundImageName string
	modify              bool
	data                *util.User
	*InfoResponse
}

func QueryUserInfo(UserId uint, SeenId uint, action_type int, nickname string, avatar *multipart.FileHeader, background_image *multipart.FileHeader, AvatarName string, BackgroundImageName string) (*InfoResponse, error) {
	return NewQueryUserInfo(UserId, SeenId, action_type, nickname, avatar, background_image, AvatarName, BackgroundImageName).Do()
}

func NewQueryUserInfo(Userid uint, SeenId uint, action_type int, nickname string, avatar *multipart.FileHeader, background_image *multipart.FileHeader, AvatarName string, BackgroundImageName string) *QueryUserInfoFlow {
	return &QueryUserInfoFlow{UserId: Userid, SeenId: SeenId, action_type: action_type, nickname: nickname, avatar: avatar, background_image: background_image, AvatarName: AvatarName, BackgroundImageName: BackgroundImageName}
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
	wg := sync.WaitGroup{}
	wg.Add(3)

	errChan := make(chan error, 3)
	defer close(errChan)

	go func() {
		defer wg.Done()
		//查action_type是否只有1或2
		if q.action_type != 1 && q.action_type != 2 {
			errStr := "action_type不是1和2，输入错误"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	//查询用户id是否在数据库中
	var UserDao = model.NewUserDao()
	go func() {
		defer wg.Done()
		if q.UserId == 0 || !UserDao.QueryUserExistByUserId(q.UserId) {
			errStr := "UserId用户不存在"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	go func() {
		defer wg.Done()
		if q.SeenId == 0 || !UserDao.QueryUserExistByUserId(q.SeenId) {
			errStr := "被看的用户不存在"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
	}

	q.modify = false
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
	//修改
	if q.modify {
		wg := sync.WaitGroup{}
		wg.Add(2)

		errChan := make(chan error, 2)
		defer close(errChan)

		go func() {
			defer wg.Done()
			if err := minio.ImageToMinio(q.avatar, q.AvatarName); err != nil {
				log.Println(err)
				errChan <- err
			}
		}()

		go func() {
			defer wg.Done()
			if err := minio.ImageToMinio(q.background_image, q.BackgroundImageName); err != nil {
				log.Println(err)
				errChan <- err
			}
		}()

		wg.Wait()
		if len(errChan) > 0 {
			return <-errChan
		}

		if err := UserDao.UpdateUserInfo(&dbUser, q.SeenId, q.nickname, q.AvatarName, q.BackgroundImageName); err != nil {
			return err
		}
	}

	if err := UserDao.QueryUserInfoById(q.SeenId, &dbUser); err != nil {
		return err
	}

	if dbUser.BackgroundImage != "" {
		dbUser.BackgroundImage = config.Miniourl + dbUser.BackgroundImage
	}
	if dbUser.Avatar != "" {
		dbUser.Avatar = config.Miniourl + dbUser.Avatar
	}

	q.data = &util.User{
		Id:              dbUser.ID,
		Name:            dbUser.UserName,
		NickName:        dbUser.NickName,
		WorkCount:       dbUser.WorkCount,
		BackGroundImage: dbUser.BackgroundImage,
		Avatar:          dbUser.Avatar,
	}
	return nil
}

func (q *QueryUserInfoFlow) packData() error {
	q.InfoResponse = &InfoResponse{MyUser: q.data}
	return nil
}
