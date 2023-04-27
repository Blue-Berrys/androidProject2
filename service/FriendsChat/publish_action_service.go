package service

import (
	"androidProject2/cache/minio"
	"androidProject2/config"
	model2 "androidProject2/model/db"
	model3 "androidProject2/model/friendschat"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"log"
	"mime/multipart"
	"sync"
)

type PublishActionResponse struct {
	FriendsChat *util.FriendsChat `json:"friendschat"`
}

type PublishActionFlow struct {
	UserId     uint
	Content    string
	ActionType int
	Images     []*multipart.FileHeader
	ImageName  []string

	ImageUrl string

	data *util.FriendsChat
	*PublishActionResponse
}

func PublishAction(userId uint, content string, action_type int, imageName []string, images []*multipart.FileHeader) (*PublishActionResponse, error) {
	return NewPublishAction(userId, content, action_type, imageName, images).Do()
}

func NewPublishAction(userId uint, context string, action_type int, imageName []string, images []*multipart.FileHeader) *PublishActionFlow {
	return &PublishActionFlow{UserId: userId, Content: context, ActionType: action_type, ImageName: imageName, Images: images}
}

func (q *PublishActionFlow) Do() (*PublishActionResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.PublishActionResponse, nil
}

func (q *PublishActionFlow) checkNum() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 3)
	defer close(errChan)

	go func() {
		defer wg.Done()
		//查action_type是否只有1或2
		if q.ActionType != 1 && q.ActionType != 2 {
			errStr := "action_type不是1和2，输入错误"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	//查数据库，这个id是否存在
	go func() {
		defer wg.Done()
		var UserDao = model.NewUserDao()
		if q.UserId == 0 || !UserDao.QueryUserExistByUserId(q.UserId) {
			errStr := "UserId用户不存在"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (q *PublishActionFlow) prepareData() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 2)
	defer close(errChan)

	var modelUser *util.User
	var dbUser = model2.User{}
	var UserDao = model.NewUserDao()
	go func() {
		defer wg.Done()
		//根据id查询用户信息
		if err := UserDao.QueryUserInfoById(q.UserId, &dbUser); err != nil {
			log.Println(err)
			errChan <- err
		}
		//构造model.user
		modelUser = &util.User{
			Id:              dbUser.ID,
			Name:            dbUser.UserName,
			NickName:        dbUser.NickName,
			WorkCount:       dbUser.WorkCount,
			BackGroundImage: dbUser.BackgroundImage,
			Level:           dbUser.Level,
			Avatar:          dbUser.Avatar,
		}
	}()

	var ResImageUrl string
	go func() {
		defer wg.Done()
		//增加朋友圈记录
		if q.ActionType == 1 { //有图片
			for i, image := range q.Images {
				if err := minio.ImageToMinio(image, q.ImageName[i]); err != nil {
					log.Println(err)
					errChan <- err
				}
				if i == 0 {
					q.ImageUrl = config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
					ResImageUrl = config.Miniourl + config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
				} else {
					q.ImageUrl += " " + config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
					ResImageUrl += " " + config.Miniourl + config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
				}
			}
		} else {
			q.ImageUrl = ""
			ResImageUrl = ""
		}
	}()
	wg.Wait()

	//插入FriendsChat
	var FriendsChatDao = model3.NewFriendsChatDao()
	var friendschat = model2.FriendsChat{
		UserId:   q.UserId,
		Content:  q.Content,
		ImageUrl: q.ImageUrl,
	}
	if err := FriendsChatDao.AddFriendsChat(&friendschat); err != nil {
		return err
	}

	//构造返回Json
	q.data = &util.FriendsChat{
		Id:         friendschat.ID,
		User:       *modelUser,
		ImageUrl:   ResImageUrl,
		Content:    q.Content,
		CreateDate: friendschat.CreatedAt.Format("01-02 15:04:05"),
	}
	//给User表中WorkCount作品数+1
	workCount := dbUser.WorkCount
	if err := UserDao.AddOneWorkCountByUserId(&dbUser, dbUser.ID, workCount+1); err != nil {
		return err
	}
	return nil
}

func (q *PublishActionFlow) packData() error {
	q.PublishActionResponse = &PublishActionResponse{FriendsChat: q.data}
	return nil
}
